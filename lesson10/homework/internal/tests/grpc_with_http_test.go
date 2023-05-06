package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	grpcSvc "homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"homework10/internal/tests/mocks"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

type ServerSuite struct {
	suite.Suite
	App     *mocks.App
	Client  *testClient
	SigQuit chan os.Signal
	eg      *errgroup.Group
}

func (suite *ServerSuite) SetupSuite() {
	appSvc := app.NewApp(adrepo.New())
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	svc := grpcSvc.NewService(appSvc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcSvc.UnaryLoggerInterceptor,
		grpcSvc.UnaryRecoveryInterceptor(),
	))
	grpcSvc.RegisterAdServiceServer(grpcServer, svc)

	httpServer := httpgin.NewHTTPServer(":18080", appSvc)
	suite.SigQuit = make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(suite.SigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		select {
		case s := <-suite.SigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})
	eg.Go(grpcSvc.RunGRPCServerGracefully(ctx, lis, grpcServer))
	eg.Go(httpgin.RunHTTPServerGracefully(ctx, httpServer))
	go func() {
		if err := eg.Wait(); err != nil {
			log.Printf("gracefully shutting down the servers: %s\n", err.Error())
		}
		log.Println("servers were successfully shutdown")
	}()
}

func (suite *ServerSuite) Test() {
	time.Sleep(6 * time.Second)
}

func (suite *ServerSuite) TearDownSuite() {
	suite.SigQuit <- syscall.SIGINT
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
