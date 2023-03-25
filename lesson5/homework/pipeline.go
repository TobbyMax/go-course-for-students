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
	// добавил проверку в начале, чтобы точно закрыть все каналы и завершились все горутины
	in = checkContext(ctx, in)
	for _, stage := range stages {
		in = stage(in)
	}
	// проверка в конце, чтобы вывод в результирующий канал прекратился сразу с отменой контекста
	out := checkContext(ctx, in)
	return out
}
