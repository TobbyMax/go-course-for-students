package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func checkContext(ctx context.Context, in In) Out {
	out := make(chan any)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()
	return out
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(out)
	}
	res := checkContext(ctx, out)
	return res
}
