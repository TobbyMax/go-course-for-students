package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/app"
)

func (s *AdService) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := s.app.CreateAd(ctx, request.Title, request.Text, request.UserId)

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}

	return AdSuccessResponse(ad), nil
}

func (s *AdService) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := s.app.ChangeAdStatus(ctx, request.AdId, request.UserId, request.Published)

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdSuccessResponse(ad), nil
}

func (s *AdService) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := s.app.UpdateAd(ctx, request.GetAdId(), request.GetUserId(), request.GetTitle(), request.GetText())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdSuccessResponse(ad), nil
}

func (s *AdService) GetAd(ctx context.Context, request *GetAdRequest) (*AdResponse, error) {
	ad, err := s.app.GetAd(ctx, request.GetAdId())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdSuccessResponse(ad), nil
}

func (s *AdService) ListAds(ctx context.Context, request *ListAdRequest) (*ListAdResponse, error) {
	date, err := ParseDate(request.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	al, err := s.app.ListAds(ctx, app.ListAdsParams{
		Published: request.Published,
		Uid:       request.UserId,
		Date:      date,
		Title:     request.Title,
	})

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return AdListSuccessResponse(al), nil
}

func (s *AdService) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	u, err := s.app.CreateUser(ctx, request.GetName(), request.GetEmail())

	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
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
		return nil, status.Error(GetErrorCode(err), err.Error())
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
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *AdService) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteAd(ctx, request.GetAdId(), request.GetAuthorId())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}