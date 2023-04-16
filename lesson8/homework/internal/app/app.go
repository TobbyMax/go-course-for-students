package app

import (
	"context"
	"fmt"
	"github.com/TobbyMax/validator"
	"homework8/internal/ads"
	"homework8/internal/user"
	"time"
)

var (
	ErrForbidden    = fmt.Errorf("forbidden")
	ErrAdNotFound   = fmt.Errorf("ad with such id does not exist")
	ErrUserNotFound = fmt.Errorf("user with such id does not exist")
)

type App interface {
	CreateAd(ctx context.Context, title string, text string, uid int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error)
	GetAd(ctx context.Context, id int64) (*ads.Ad, error)

	ListAds(ctx context.Context, params ListAdsParams) (*ads.AdList, error)

	CreateUser(ctx context.Context, nickname string, email string) (*user.User, error)
	GetUser(ctx context.Context, id int64) (*user.User, error)
	UpdateUser(ctx context.Context, id int64, nickname string, email string) (*user.User, error)
}

type AdRepository interface {
	AddAd(ctx context.Context, ad ads.Ad) (int64, error)
	GetAdByID(ctx context.Context, id int64) (*ads.Ad, error)
	UpdateAdStatus(ctx context.Context, id int64, published bool, date time.Time) error
	UpdateAdContent(ctx context.Context, id int64, title string, text string, date time.Time) error

	GetAdList(ctx context.Context, params ListAdsParams) (*ads.AdList, error)
}

type UserRepository interface {
	AddUser(ctx context.Context, u user.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
	UpdateUser(ctx context.Context, id int64, nickname string, email string) error
}

type Repository interface {
	AdRepository
	UserRepository
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
	ad := ads.Ad{Title: title, Text: text, AuthorID: uid, Published: false, DateCreated: time.Now().UTC()}
	ad.DateChanged = ad.DateCreated
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
	ad.DateChanged = time.Now().UTC()

	err = a.repository.UpdateAdStatus(ctx, id, published, ad.DateChanged)
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
	ad.DateChanged = time.Now().UTC()

	if err := validator.Validate(*ad); err != nil {
		return nil, err
	}

	err = a.repository.UpdateAdContent(ctx, id, title, text, ad.DateChanged)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a AdApp) ListAds(ctx context.Context, params ListAdsParams) (*ads.AdList, error) {
	p := true
	if params.Published == nil && params.Uid == nil && params.Date == nil && params.Title == nil {
		params.Published = &p
	}
	al, err := a.repository.GetAdList(ctx, params)

	if err != nil {
		return nil, err
	}

	return al, nil
}

func (a AdApp) CreateUser(ctx context.Context, nickname string, email string) (*user.User, error) {
	u := user.User{Nickname: nickname, Email: email}

	if err := validator.Validate(u); err != nil {
		return nil, err
	}

	id, err := a.repository.AddUser(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = id

	return &u, nil
}

func (a AdApp) GetUser(ctx context.Context, id int64) (*user.User, error) {
	u, err := a.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a AdApp) UpdateUser(ctx context.Context, id int64, nickname string, email string) (*user.User, error) {
	u, err := a.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Nickname = nickname
	u.Email = email

	if err := validator.Validate(*u); err != nil {
		return nil, err
	}

	err = a.repository.UpdateUser(ctx, id, nickname, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
