package app

import (
	"context"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx context.Context, title string, text string, uid int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error)
}

type Repository interface {
	AddAd(ctx context.Context, ad ads.Ad) (int64, error)
	ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error)
	ChangeAdContent(id int64, uid int64, title string, text string) (*ads.Ad, error)
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

	id, err := a.repository.AddAd(ctx, ad)
	if err != nil {
		return nil, err
	}
	ad.ID = id

	return &ad, nil
}

func (a AdApp) ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error) {
	ad, err := a.repository.ChangeAdStatus(ctx, id, uid, published)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a AdApp) UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error) {
	ad, err := a.repository.ChangeAdContent(id, uid, title, text)
	if err != nil {
		return nil, err
	}

	return ad, nil
}
