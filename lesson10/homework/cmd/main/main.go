package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	"homework10/internal/graceful"
	grpcSvc "homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"os"

	"log"
	"net"
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
	eg.Go(graceful.CaptureSignal(ctx, sigQuit))
	// run grpc server
	eg.Go(grpcSvc.RunGRPCServerGracefully(ctx, lis, grpcServer))
	// run http server
	eg.Go(httpgin.RunHTTPServerGracefully(ctx, httpServer))

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}
	log.Println("servers were successfully shutdown")
}
