package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("TobbyMax", "agemax@gmail.com")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, "TobbyMax", response.Data.Nickname)
	assert.Equal(t, "agemax@gmail.com", response.Data.Email)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("TobbyMax", "abc")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("TobbyMax", "agemax@gmail.com")
	assert.NoError(t, err)

	response, err = client.getUser(response.Data.ID)
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, "TobbyMax", response.Data.Nickname)
	assert.Equal(t, "agemax@gmail.com", response.Data.Email)
}

func TestGetUser_NonExistentID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("MacMiller", "blue_slide_park@hotmail.com")
	assert.NoError(t, err)

	_, err = client.getUser(1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("MacMiller", "swimming@circles.com")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.ID, "MacMiller", "the_divine2016@feminine.ru")
	assert.NoError(t, err)
	assert.Equal(t, "MacMiller", response.Data.Nickname)
	assert.Equal(t, "the_divine2016@feminine.ru", response.Data.Email)
}

func TestUpdateUser_InvalidEmail(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("MacMiller", "swimming@circles.com")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.ID, "MacMiller", "good_am.ru")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_ID(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("Mac Miller", "swimming@circles.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createUser("Mac Miller", "swimming@circles.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createUser("Mac Miller", "swimming@circles.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("Mac Miller", "swimming@circles.com")
	assert.NoError(t, err)

	ad1, err := client.createAd(user1.Data.ID, "Good News", "Dang!")
	assert.NoError(t, err)

	_, err = client.createUser("Mac Miller", "swimming@circles.com")
	assert.NoError(t, err)

	_, err = client.deleteUser(user1.Data.ID)
	assert.NoError(t, err)

	_, err = client.getUser(user1.Data.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, ErrNotFound, err)

	_, err = client.getAd(ad1.Data.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}
