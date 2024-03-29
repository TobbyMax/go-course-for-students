package graceful

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CaptureSignal(ctx context.Context) func() error {
	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	return func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	}
}
