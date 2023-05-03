package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(uResponse.Data.ID))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(uResponse.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(uResponse.Data.ID, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(uResponse.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(uResponse.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestGetAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.getAd(response.Data.ID)
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(uResponse.Data.ID))
	assert.False(t, response.Data.Published)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("Mac Miller", "swimming@circles.com")
	assert.NoError(t, err)

	ad1, err := client.createAd(user1.Data.ID, "Good News", "Dang!")
	assert.NoError(t, err)

	_, err = client.deleteAd(ad1.Data.ID, user1.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(ad1.Data.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}
