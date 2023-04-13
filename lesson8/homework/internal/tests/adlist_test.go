package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestListAdsPublished(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "best cat", "not for sale")
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

func TestListAdsNotPublished(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	notPublishedAd, err := client.createAd(user2.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByStatus(false)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, notPublishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, notPublishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, notPublishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, notPublishedAd.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestListAdsByUser(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	adByUser1, err := client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByUser(user1.Data.ID)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, adByUser1.Data.ID)
	assert.Equal(t, ads.Data[0].Title, adByUser1.Data.Title)
	assert.Equal(t, ads.Data[0].Text, adByUser1.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, adByUser1.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestListAdsByDate(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByDate(time.Now())
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 3)
}

func TestListAdsByDate_Yesterday(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByDate(time.Now().Add(time.Duration(-24) * time.Hour))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 0)
}

func TestListAdsByUserAndDate(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByUserAndDate(user2.Data.ID, time.Now())
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

func TestListAdsByTitle(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	titledAd, err := client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByTitle(titledAd.Data.Title)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, titledAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, titledAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, titledAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, titledAd.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestListAdsByTitle_Multiple(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	titledAd, err := client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	response, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "GOMD", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsByTitle(titledAd.Data.Title)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

func TestListAdsByOptions(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	user2, err := client.createUser("Kendrick", "section80@damn.com")
	assert.NoError(t, err)

	response, err := client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	assert.NoError(t, err)

	target, err := client.changeAdStatus(user1.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "GOMD", "Cole World")
	assert.NoError(t, err)

	response, err = client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	response, err = client.createAd(user2.Data.ID, "GOMD", "not for sale")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.listAdsByOptions(user1.Data.ID, time.Now(), true, target.Data.Title)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, target.Data.ID)
	assert.Equal(t, ads.Data[0].Title, target.Data.Title)
	assert.Equal(t, ads.Data[0].Text, target.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, target.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}
