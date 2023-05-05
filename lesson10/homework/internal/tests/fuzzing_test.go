package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func FuzzCreateUserID_Fuzz(f *testing.F) {
	client := getTestClient()
	for i := 0; i < 100; i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, id int) {
		name := "Mac Miller"
		email := "swimming@circles.com"
		resp, err := client.createUser(name, email)

		assert.NoError(t, err)
		assert.Equal(t, int64(id), resp.Data.ID)
		assert.Equal(t, name, resp.Data.Nickname)
		assert.Equal(t, email, resp.Data.Email)
	})
}

func FuzzGetUserID_Fuzz(f *testing.F) {
	client := getTestClient()
	for i := 0; i < 100; i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, id int) {
		name := "Mac Miller"
		email := "swimming@circles.com"
		_, err := client.createUser(name, email)
		assert.NoError(t, err)

		resp, err := client.getUser(int64(id))
		assert.NoError(t, err)

		assert.Equal(t, int64(id), resp.Data.ID)
		assert.Equal(t, name, resp.Data.Nickname)
		assert.Equal(t, email, resp.Data.Email)
	})
}
