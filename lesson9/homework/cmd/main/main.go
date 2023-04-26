package main

import (
	"google.golang.org/grpc"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcSvc "homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"
	"log"
	"net"
	"sync"
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
	grpcServer := grpc.NewServer()
	grpcSvc.RegisterAdServiceServer(grpcServer, svc)

	server := httpgin.NewHTTPServer(httpPort, appSvc)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		defer wg.Done()
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	wg.Wait()
}
