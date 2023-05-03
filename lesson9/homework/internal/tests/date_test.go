package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateDate(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, uResponse.Data.ID)
	assert.False(t, response.Data.Published)

	assert.True(t, response.Data.DateCreated == response.Data.DateChanged)
	date, _ := time.Parse(DateTimeLayout, response.Data.DateCreated)
	assert.True(t, time.Since(date) < time.Hour)
}

func TestChangeDate(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("J.Cole", "foresthill@drive.com")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	response, err = client.changeAdStatus(uResponse.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
	
	assert.True(t, response.Data.DateCreated != response.Data.DateChanged)
	date, _ := time.Parse(DateTimeLayout, response.Data.DateChanged)
	assert.True(t, time.Since(date) < time.Hour)
}
