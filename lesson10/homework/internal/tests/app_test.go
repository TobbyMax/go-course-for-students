package tests

import (
	"context"
	"github.com/TobbyMax/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/tests/mocks"
	"homework10/internal/user"
	"testing"
)

func TestApp_CreateAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddAd", ctx, mock.AnythingOfType("ads.Ad")).
		Return(id, nil).
		Once()
	service := app.NewApp(repo)
	ad, err := service.CreateAd(ctx, "title", "text", 1)
	assert.Nil(t, err)
	assert.Equal(t, id, ad.ID)
	assert.Equal(t, "title", ad.Title)
	assert.Equal(t, "text", ad.Text)
	assert.Equal(t, int64(1), ad.AuthorID)
}

func TestApp_CreateAd_NonExistentUser(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddAd", ctx, mock.AnythingOfType("ads.Ad")).
		Return(id, app.ErrUserNotFound).
		Once()
	service := app.NewApp(repo)
	_, err := service.CreateAd(ctx, "title", "text", 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrUserNotFound)
}

func TestApp_GetAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	repo.On("GetAdByID", ctx, int64(0)).
		Return(&ads.Ad{}, nil).
		Once()
	service := app.NewApp(repo)
	ad, err := service.GetAd(ctx, 0)
	assert.Nil(t, err)
	assert.Equal(t, ads.Ad{}, *ad)
}

func TestApp_CreateAd_InvalidTitle(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddAd", ctx, mock.AnythingOfType("ads.Ad")).
		Return(id, nil).
		Once()
	service := app.NewApp(repo)
	_, err := service.CreateAd(ctx, "", "text", 1)
	assert.Error(t, err)
	e := &validator.ValidationErrors{}
	assert.ErrorAs(t, err, e)
}

func TestApp_CreateAd_InvalidText(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddAd", ctx, mock.AnythingOfType("ads.Ad")).
		Return(id, nil).
		Once()
	service := app.NewApp(repo)
	_, err := service.CreateAd(ctx, "title", "", 1)
	assert.Error(t, err)
	e := &validator.ValidationErrors{}
	assert.ErrorAs(t, err, e)
}

func TestApp_UpdateAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	title := "title"
	text := "text"
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	repo.On("UpdateAdContent", ctx, id, title, text, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateAd(ctx, id, int64(1), title, text)
	assert.Nil(t, err)
}

func TestApp_UpdateAd_NonExistentAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	title := "title"
	text := "text"
	repo.On("GetAdByID", ctx, id).
		Return(nil, app.ErrAdNotFound).
		Once()
	repo.On("UpdateAdContent", ctx, id, title, text, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateAd(ctx, id, int64(1), title, text)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrAdNotFound)
}

func TestApp_UpdateAd_Forbidden(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	title := "title"
	text := "text"
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 0}, nil).
		Once()
	repo.On("UpdateAdContent", ctx, id, title, text, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateAd(ctx, id, int64(1), title, text)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrForbidden)
}

func TestApp_UpdateAd_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	title := "title"
	text := "text"
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	repo.On("UpdateAdContent", ctx, id, title, text, mock.AnythingOfType("time.Time")).
		Return(ErrMock).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateAd(ctx, id, int64(1), title, text)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMock)
}

func TestApp_ChangeAdStatus(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	repo.On("UpdateAdStatus", ctx, id, true, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.ChangeAdStatus(ctx, id, int64(1), true)
	assert.Nil(t, err)
}

func TestApp_ChangeAdStatus_NonExistentAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(nil, app.ErrAdNotFound).
		Once()
	repo.On("UpdateAdStatus", ctx, id, true, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.ChangeAdStatus(ctx, id, int64(1), true)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrAdNotFound)
}

func TestApp_ChangeAdStatus_Forbidden(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 0}, nil).
		Once()
	repo.On("UpdateAdStatus", ctx, id, true, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.ChangeAdStatus(ctx, id, int64(1), true)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrForbidden)
}

func TestApp_ChangeAdStatus_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	repo.On("UpdateAdStatus", ctx, id, true, mock.AnythingOfType("time.Time")).
		Return(ErrMock).
		Once()

	service := app.NewApp(repo)
	_, err := service.ChangeAdStatus(ctx, id, int64(1), true)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMock)
}

func TestApp_ListAds(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	pub := true
	params := app.ListAdsParams{Published: &pub}
	repo.On("GetAdList", ctx, params).
		Return(nil, nil).
		Once()

	service := app.NewApp(repo)
	al, err := service.ListAds(ctx, params)
	assert.Nil(t, err)
	assert.Nil(t, al)
}

func TestApp_ListAds_AllNil(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	params := app.ListAdsParams{}
	pub := true
	repo.On("GetAdList", ctx, app.ListAdsParams{Published: &pub}).
		Return(nil, nil).
		Once()

	service := app.NewApp(repo)
	al, err := service.ListAds(ctx, params)
	assert.Nil(t, err)
	assert.Nil(t, al)
}

func TestApp_ListAds_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	params := app.ListAdsParams{}
	pub := true
	repo.On("GetAdList", ctx, app.ListAdsParams{Published: &pub}).
		Return(nil, ErrMock).
		Once()

	service := app.NewApp(repo)
	al, err := service.ListAds(ctx, params)
	assert.Nil(t, al)
	assert.ErrorIs(t, err, ErrMock)
}

func TestApp_CreateUser(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddUser", ctx, mock.AnythingOfType("user.User")).
		Return(id, nil).
		Once()
	service := app.NewApp(repo)
	u, err := service.CreateUser(ctx, "Mac Miller", "swimming@circles.com")
	assert.Nil(t, err)
	assert.Equal(t, id, u.ID)
	assert.Equal(t, "Mac Miller", u.Nickname)
	assert.Equal(t, "swimming@circles.com", u.Email)
}

func TestApp_GetUser(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	repo.On("GetUserByID", ctx, int64(0)).
		Return(nil, nil).
		Once()
	service := app.NewApp(repo)
	u, err := service.GetUser(ctx, 0)
	assert.Nil(t, err)
	assert.Nil(t, u)
}

func TestApp_CreateUser_InvalidName(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddUser", ctx, mock.AnythingOfType("user.User")).
		Return(id, nil).
		Once()
	service := app.NewApp(repo)
	u, err := service.CreateUser(ctx, "", "swimming@circles.com")
	assert.Error(t, err)
	e := &validator.ValidationErrors{}
	assert.ErrorAs(t, err, e)
	assert.Nil(t, u)
}

func TestApp_CreateUser_InvalidEmail(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddUser", ctx, mock.AnythingOfType("user.User")).
		Return(id, nil).
		Once()
	service := app.NewApp(repo)
	u, err := service.CreateUser(ctx, "Mac", "")
	assert.Error(t, err)
	e := &validator.ValidationErrors{}
	assert.ErrorAs(t, err, e)
	assert.Nil(t, u)
}

func TestApp_CreateUser_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(13)
	repo.On("AddUser", ctx, mock.AnythingOfType("user.User")).
		Return(id, ErrMock).
		Once()
	service := app.NewApp(repo)
	u, err := service.CreateUser(ctx, "Mac", "swimming@circles.com")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMock)
	assert.Nil(t, u)
}

func TestApp_UpdateUser(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	name := "Mac Miller"
	email := "swimming@circles.com"
	repo.On("GetUserByID", ctx, id).
		Return(&user.User{}, nil).
		Once()
	repo.On("UpdateUser", ctx, id, name, email).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateUser(ctx, id, name, email)
	assert.Nil(t, err)
}

func TestApp_UpdateUser_NonExistentID(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	name := "Mac Miller"
	email := "swimming@circles.com"
	repo.On("GetUserByID", ctx, id).
		Return(nil, app.ErrUserNotFound).
		Once()
	repo.On("UpdateUser", ctx, id, name, email).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateUser(ctx, id, name, email)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrUserNotFound)
}

func TestApp_UpdateUser_InvalidName(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	name := ""
	email := "swimming@circles.com"
	repo.On("GetUserByID", ctx, id).
		Return(&user.User{}, nil).
		Once()
	repo.On("UpdateUser", ctx, id, name, email).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateUser(ctx, id, name, email)
	assert.Error(t, err)
	e := &validator.ValidationErrors{}
	assert.ErrorAs(t, err, e)
}

func TestApp_UpdateUser_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	name := "Mac Miller"
	email := "swimming@circles.com"
	repo.On("GetUserByID", ctx, id).
		Return(&user.User{}, nil).
		Once()
	repo.On("UpdateUser", ctx, id, name, email).
		Return(ErrMock).
		Once()

	service := app.NewApp(repo)
	_, err := service.UpdateUser(ctx, id, name, email)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMock)
}

func TestApp_DeleteUser(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("DeleteUserByID", ctx, id).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	err := service.DeleteUser(ctx, id)
	assert.Nil(t, err)
}

func TestApp_DeleteUser_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("DeleteUserByID", ctx, id).
		Return(ErrMock).
		Once()

	service := app.NewApp(repo)
	err := service.DeleteUser(ctx, id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMock)
}

func TestApp_DeleteAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	repo.On("DeleteAdByID", ctx, id).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	err := service.DeleteAd(ctx, id, 1)
	assert.Nil(t, err)
}

func TestApp_DeleteAd_Forbidden(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 0}, nil).
		Once()
	repo.On("DeleteAdByID", ctx, id).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	err := service.DeleteAd(ctx, id, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrForbidden)
}

func TestApp_DeleteAd_NonExistentAd(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(nil, app.ErrAdNotFound).
		Once()
	repo.On("DeleteAdByID", ctx, id).
		Return(nil).
		Once()

	service := app.NewApp(repo)
	err := service.DeleteAd(ctx, id, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrAdNotFound)
}

func TestApp_DeleteAd_RepoError(t *testing.T) {
	repo := &mocks.Repository{}
	ctx := context.Background()
	id := int64(0)
	repo.On("GetAdByID", ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	repo.On("DeleteAdByID", ctx, id).
		Return(ErrMock).
		Once()

	service := app.NewApp(repo)
	err := service.DeleteAd(ctx, id, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMock)
}
