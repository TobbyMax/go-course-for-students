package grpc

import (
	"context"
	"google.golang.org/grpc"
	"homework9/internal/app"
	"log"
	"time"
)

type AdService struct {
	app app.App
}

func NewService(a app.App) AdServiceServer {
	service := &AdService{app: a}
	return service
}

func LoggerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	log.Printf("-- received request -- | protocol: GRPC | method: %s", info.FullMethod)

	h, err := handler(ctx, req)

	latency := time.Since(start)
	log.Printf("-- handled request -- | protocol: GRPC | latency: %+v | method: %s | error: (%v)\n",
		latency, info.FullMethod, err)

	return h, err
}
