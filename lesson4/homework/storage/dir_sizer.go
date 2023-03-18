package storage

import (
	"context"
	"runtime"
	"sync"
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

	// TODO: add other fields as you wish
	subDirs chan Dir
	result  chan Result
}

const WorkerCount = 2

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: WorkerCount,
		subDirs:         make(chan Dir, 10000),
		result:          make(chan Result, WorkerCount+1),
	}
}

func worker2(ctx context.Context, wg *sync.WaitGroup, dirs <-chan Dir, result chan Result, errs chan<- error) {
	defer wg.Done()
	workerRes := Result{}
	for {
		select {
		case dir, ok := <-dirs:
			if !ok {
				res := <-result
				res.Size += workerRes.Size
				res.Count += workerRes.Count
				result <- res
				return
			}
			_, files, err := dir.Ls(ctx)
			if err != nil {
				//errs <- err
				return
			}
			for _, file := range files {
				fileSize, err := file.Stat(ctx)
				if err != nil {
					//errs <- err
					return
				}
				runtime.Gosched()
				workerRes.Count++
				workerRes.Size += fileSize
			}
		case <-ctx.Done():
			res := <-result
			res.Size += workerRes.Size
			res.Count += workerRes.Count
			result <- res
			return
		}
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, dirs chan Dir, result chan Result, errs chan<- error) {
	defer wg.Done()
	//workerRes := Result{}
	for len(dirs) > 0 {
		select {
		case dir, ok := <-dirs:
			if !ok {
				//result <- workerRes
				return
			}
			subDirs, files, err := dir.Ls(ctx)
			if err != nil {
				errs <- err
				return
			}
			for _, d := range subDirs {
				dirs <- d
			}
			workerRes := Result{}
			for _, file := range files {
				fileSize, err := file.Stat(ctx)
				if err != nil {
					errs <- err
					return
				}
				workerRes.Size += fileSize
			}
			runtime.Gosched()
			res := <-result
			res.Size += workerRes.Size
			res.Count += int64(len(files))
			result <- res
		case <-ctx.Done():
			//result <- workerRes
			return
		}
	}
	errs <- nil
}

func (s *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup
	ctxWithCancel, cancel := context.WithCancel(ctx)
	s.result <- Result{}
	errs := make(chan error)
	defer cancel()

	//res := Result{}
	s.subDirs <- d
	for i := 0; i < s.maxWorkersCount; i++ {
		wg.Add(1)
		go worker(ctxWithCancel, &wg, s.subDirs, s.result, errs)
	}
	//wg.Wait()
	//res = <-s.result
	//return res, nil\
	count := 0
	for {
		select {
		case err, _ := <-errs:
			if err == nil {
				count++
			} else {
				return Result{}, err
			}
			if count == s.maxWorkersCount {
				cancel()
				return <-s.result, nil
			}
		}
	}
}
