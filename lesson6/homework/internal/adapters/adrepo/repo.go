package adrepo

import (
	"context"
	"fmt"
	"homework6/internal/ads"
	"homework6/internal/app"
)

var ErrIDNotFound = fmt.Errorf("ad with such id does not exist")
var ErrForbidden = fmt.Errorf("forbidden")

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

func (r RepositoryMap) ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error) {
	ad := &ads.Ad{}
	if _, ok := r.data[id]; !ok {
		return ad, ErrIDNotFound
	}
	if r.data[id].AuthorID != uid {
		return ad, ErrForbidden
	}
	*ad = r.data[id]
	ad.Published = published
	return ad, nil
}

func (r RepositoryMap) ChangeAdContent(id int64, uid int64, title string, text string) (*ads.Ad, error) {
	ad := &ads.Ad{}
	if _, ok := r.data[id]; !ok {
		return ad, ErrIDNotFound
	}
	if r.data[id].AuthorID != uid {
		return ad, ErrForbidden
	}
	*ad = r.data[id]
	ad.Title = title
	ad.Text = text
	return ad, nil
}
