package httpfiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"homework6/internal/ads"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func (r createAdRequest) validate() error {
	if r.Title == "" {
		return errors.New("empty title")
	}
	if r.Text == "" {
		return errors.New("empty text")
	}
	if len(r.Title) >= 100 {
		return errors.New("empty title")
	}
	if len(r.Text) >= 500 {
		return errors.New("empty text")
	}
	return nil
}

type adResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func (r updateAdRequest) validate() error {
	if r.Title == "" {
		return errors.New("empty title")
	}
	if r.Text == "" {
		return errors.New("empty text")
	}
	if len(r.Title) >= 100 {
		return errors.New("empty title")
	}
	if len(r.Text) >= 500 {
		return errors.New("empty text")
	}
	return nil
}

func AdSuccessResponse(ad *ads.Ad) *fiber.Map {
	return &fiber.Map{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		},
		"error": nil,
	}
}

func AdErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"data":  nil,
		"error": err.Error(),
	}
}
