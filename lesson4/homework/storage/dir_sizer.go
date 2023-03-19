package storage

import (
	"context"
	"math"
	"runtime"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
	// dirs queue of subdirectories to calculate size
	dirs chan Dir
	// queueLen length of dirs queue
	queueLen int64
}

const WorkerCount = 4

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: WorkerCount,
	}
}

// worker a unit that calculates Size
func (s *sizer) worker(ctx context.Context, result *Result, errs chan<- error) {
	for {
		select {
		case dir, ok := <-s.dirs:
			if !ok {
				return
			}
			subDirs, files, err := dir.Ls(ctx)
			if err != nil {
				errs <- err
				return
			}
			for _, d := range subDirs {
				s.dirs <- d
			}

			dirSize := int64(0)
			for _, file := range files {
				fileSize, err := file.Stat(ctx)
				if err != nil {
					errs <- err
					return
				}
				dirSize += fileSize
			}
			atomic.AddInt64(&result.Count, int64(len(files)))
			atomic.AddInt64(&result.Size, dirSize)
			atomic.AddInt64(&s.queueLen, int64(len(subDirs))-1)

		case <-ctx.Done():
			return
		}
	}
}

// Size calculates size of a given directory
func (s *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	runtime.GOMAXPROCS(4)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	result := Result{}
	errs := make(chan error)
	defer close(errs)

	s.dirs = make(chan Dir, math.MaxInt16)
	defer close(s.dirs)
	s.dirs <- d
	s.queueLen++

	for i := 0; i < s.maxWorkersCount; i++ {
		go s.worker(ctxWithCancel, &result, errs)
	}

	for {
		select {
		case err := <-errs:
			return Result{}, err
		default:
			if s.queueLen == 0 {
				return result, nil
			}
		}
	}
}
