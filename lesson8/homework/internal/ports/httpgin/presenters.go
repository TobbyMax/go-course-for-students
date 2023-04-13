package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	"homework8/internal/user"
	"time"
)

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	AuthorID    int64     `json:"author_id"`
	Published   bool      `json:"published"`
	DateCreated time.Time `json:"date_created"`
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

type listAdsRequest struct {
	Published *bool      `json:"published"`
	UserID    *int64     `json:"user_id"`
	Date      *time.Time `json:"date"`
}

type adListResponse []adResponse

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:          ad.ID,
			Title:       ad.Title,
			Text:        ad.Text,
			AuthorID:    ad.AuthorID,
			Published:   ad.Published,
			DateCreated: ad.DateCreated,
		},
		"error": nil,
	}
}

func AdListSuccessResponse(al *ads.AdList) *gin.H {
	data := make(adListResponse, 0)
	for _, ad := range al.Data {
		data = append(data,
			adResponse{
				ID:          ad.ID,
				Title:       ad.Title,
				Text:        ad.Text,
				AuthorID:    ad.AuthorID,
				Published:   ad.Published,
				DateCreated: ad.DateCreated,
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

func UserSuccessResponse(u *user.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       u.ID,
			Nickname: u.Nickname,
			Email:    u.Email,
		},
		"error": nil,
	}
}

func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
