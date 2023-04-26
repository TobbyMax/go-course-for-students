package grpc

import (
	"errors"
	"github.com/TobbyMax/validator"
	"google.golang.org/grpc/codes"
	"homework9/internal/ads"
	"homework9/internal/app"
	"time"
)

func AdSuccessResponse(ad *ads.Ad) *AdResponse {
	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}
}

func AdListSuccessResponse(al *ads.AdList) *ListAdResponse {
	response := ListAdResponse{List: make([]*AdResponse, 0)}

	for _, ad := range al.Data {
		response.List = append(response.List, AdSuccessResponse(&ad))
	}
	return &response
}

func GetErrorCode(err error) codes.Code {
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

func ParseDate(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	date, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &date, nil
}
