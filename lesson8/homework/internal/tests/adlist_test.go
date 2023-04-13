package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListAdsPublished(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
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
