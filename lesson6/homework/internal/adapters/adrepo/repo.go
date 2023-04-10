package adrepo

import (
	"context"
	"fmt"
	"homework6/internal/ads"
	"homework6/internal/app"
)

var ErrIDNotFound = fmt.Errorf("ad with such id does not exist")

func New() app.Repository {
	return NewRepositorySlice()
}

type RepositoryMap struct {
	data map[int64]ads.Ad
}

func NewRepositorySlice() *RepositoryMap {
	return &RepositoryMap{data: make(map[int64]ads.Ad)}
}

func (r RepositoryMap) AddAd(ctx context.Context, ad ads.Ad) (int64, error) {
	ad.ID = int64(len(r.data))
	r.data[ad.ID] = ad
	return ad.ID, nil
}

func (r RepositoryMap) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {
	if ad, ok := r.data[id]; !ok {
		return nil, ErrIDNotFound
	} else {
		return &ad, nil
	}
}

func (r RepositoryMap) UpdateAdStatus(ctx context.Context, id int64, published bool) error {
	if _, ok := r.data[id]; !ok {
		return ErrIDNotFound
	}
	ad := r.data[id]
	ad.Published = published
	r.data[id] = ad
	return nil
}

func (r RepositoryMap) UpdateAdContent(ctx context.Context, id int64, title string, text string) error {
	if _, ok := r.data[id]; !ok {
		return ErrIDNotFound
	}
	ad := r.data[id]
	ad.Title = title
	ad.Text = text
	r.data[id] = ad
	return nil
}
