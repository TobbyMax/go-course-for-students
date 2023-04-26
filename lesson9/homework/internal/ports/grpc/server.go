package grpc

import (
	"homework9/internal/app"
)

type AdService struct {
	app app.App
}

func NewService(a app.App) AdServiceServer {
	service := &AdService{
		app: a,
	}

	return service
}
