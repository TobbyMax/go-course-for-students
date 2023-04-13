package app

import (
	"context"
	"fmt"
	"github.com/TobbyMax/validator"
	"homework8/internal/ads"
)

var ErrForbidden = fmt.Errorf("forbidden")

type App interface {
	CreateAd(ctx context.Context, title string, text string, uid int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error)
	GetAd(ctx context.Context, id int64) (*ads.Ad, error)
	ListAds(ctx context.Context) (*ads.AdList, error)

	//CreateUser(ctx context.Context, nickname string, email string) (*user.User, error)
	//GetUser(ctx context.Context, id int64) (*user.User, error)
	//UpdateUser(ctx context.Context, id int64, nickname string, email string) (*user.User, error)
}

type Repository interface {
	AddAd(ctx context.Context, ad ads.Ad) (int64, error)
	GetAdByID(ctx context.Context, id int64) (*ads.Ad, error)
	UpdateAdStatus(ctx context.Context, id int64, published bool) error
	UpdateAdContent(ctx context.Context, id int64, title string, text string) error
	GetPublishedAds(ctx context.Context) (*ads.AdList, error)

	//AddUser(ctx context.Context, u user.User) (int64, error)
	//GetUserByID(ctx context.Context, id int64) (*user.User, error)
	//UpdateUser(ctx context.Context, id int64, nickname string, email string) error
}

type AdApp struct {
	repository Repository
}

func NewApp(repo Repository) App {
	return NewAdApp(repo)
}

func NewAdApp(repo Repository) *AdApp {
	return &AdApp{repository: repo}
}

func (a AdApp) CreateAd(ctx context.Context, title string, text string, uid int64) (*ads.Ad, error) {
	ad := ads.Ad{Title: title, Text: text, AuthorID: uid, Published: false}

	if err := validator.Validate(ad); err != nil {
		return nil, err
	}

	id, err := a.repository.AddAd(ctx, ad)
	if err != nil {
		return nil, err
	}
	ad.ID = id

	return &ad, nil
}

func (a AdApp) GetAd(ctx context.Context, id int64) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a AdApp) ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != uid {
		return nil, ErrForbidden
	}

	ad.Published = published

	err = a.repository.UpdateAdStatus(ctx, id, published)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a AdApp) UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != uid {
		return nil, ErrForbidden
	}

	ad.Title = title
	ad.Text = text

	if err := validator.Validate(*ad); err != nil {
		return nil, err
	}

	err = a.repository.UpdateAdContent(ctx, id, title, text)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a AdApp) ListAds(ctx context.Context) (*ads.AdList, error) {
	al, err := a.repository.GetPublishedAds(ctx)

	if err != nil {
		return nil, err
	}

	return al, nil
}

//func (a AdApp) CreateUser(ctx context.Context, nickname string, email string) (*user.User, error) {
//	u := user.User{Nickname: nickname, Email: email}
//
//	if err := validator.Validate(u); err != nil {
//		return nil, err
//	}
//
//	id, err := a.repository.AddUser(ctx, u)
//	if err != nil {
//		return nil, err
//	}
//	u.ID = id
//
//	return &u, nil
//}
