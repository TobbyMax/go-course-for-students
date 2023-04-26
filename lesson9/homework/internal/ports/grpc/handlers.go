package grpc

import (
	"context"
	"errors"
	"github.com/TobbyMax/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/app"
)

func (s *AdService) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := s.app.CreateAd(ctx, request.Title, request.Text, request.UserId)

	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}

	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}, nil
}

func (s *AdService) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := s.app.ChangeAdStatus(ctx, request.AdId, request.UserId, request.Published)

	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}
	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}, nil
}

func (s *AdService) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := s.app.UpdateAd(ctx, request.GetAdId(), request.GetUserId(), request.GetTitle(), request.GetText())

	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}
	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}, nil
}

func (s *AdService) ListAds(ctx context.Context, empty *emptypb.Empty) (*ListAdResponse, error) {
	panic("GO")
}

func (s *AdService) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	u, err := s.app.CreateUser(ctx, request.GetName(), request.GetEmail())

	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}
	return &UserResponse{
		Id:    u.ID,
		Name:  u.Nickname,
		Email: u.Email,
	}, nil
}

func (s *AdService) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	u, err := s.app.GetUser(ctx, request.GetId())

	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}
	return &UserResponse{
		Id:    u.ID,
		Name:  u.Nickname,
		Email: u.Email,
	}, nil
}

func (s *AdService) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *AdService) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteAd(ctx, request.GetAdId(), request.GetAuthorId())
	if err != nil {
		return nil, status.Error(s.getErrorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *AdService) getErrorCode(err error) codes.Code {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		return codes.InvalidArgument
	case errors.Is(err, app.ErrForbidden):
		return codes.PermissionDenied
	case errors.Is(err, app.ErrAdNotFound):
		fallthrough
	case errors.Is(err, app.ErrUserNotFound):
		return codes.NotFound
	}
	return codes.Internal
}
