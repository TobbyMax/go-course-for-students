package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
)

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
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

type adListResponse []adResponse

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
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

func AdListSuccessResponse(al *ads.AdList) *gin.H {
	data := make(adListResponse, 0)
	for _, ad := range al.Data {
		data = append(data,
			adResponse{
				ID:        ad.ID,
				Title:     ad.Title,
				Text:      ad.Text,
				AuthorID:  ad.AuthorID,
				Published: ad.Published,
			})
	}
	return &gin.H{
		"data":  data,
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
