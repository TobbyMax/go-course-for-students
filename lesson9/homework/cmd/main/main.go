package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcSvc "homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"

	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcPort = ":8080"
	httpPort = ":18080"
)

func main() {
	appSvc := app.NewApp(adrepo.New())

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := grpcSvc.NewService(appSvc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcSvc.UnaryLoggerInterceptor,
		grpcSvc.UnaryRecoveryInterceptor(),
	))
	grpcSvc.RegisterAdServiceServer(grpcServer, svc)

	httpServer := httpgin.NewHTTPServer(httpPort, appSvc)

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(func() error {
		log.Printf("starting grpc server, listening on %s\n", grpcPort)
		defer log.Printf("close grpc server listening on %s\n", grpcPort)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	// run http server
	eg.Go(func() error {
		log.Printf("starting http server, listening on %s\n", httpServer.Addr)
		defer log.Printf("close http server listening on %s\n", httpServer.Addr)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", httpServer.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
