package tests

import (
	"github.com/TobbyMax/validator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/app"
	"homework10/internal/ports/httpgin"
	"homework10/internal/tests/mocks"
	"homework10/internal/user"
	"net/http/httptest"
	"testing"
)

type HTTPSuite struct {
	suite.Suite
	App    *mocks.App
	Client *testClient
}

func (suite *HTTPSuite) SetupTest() {
	suite.App = &mocks.App{}
	server := httpgin.NewHTTPServer(":18080", suite.App)
	testServer := httptest.NewServer(server.Handler)

	suite.Client = &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (suite *HTTPSuite) TestHandler_CreateUser() {
	type args struct {
		badReq   bool
		nickname string
		email    string
		err      error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful create",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			wantErr: false,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "invalid email",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming.com",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "validation error",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      validator.ValidationErrors{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      ErrMock,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.App.On("CreateUser",
				mock.AnythingOfType("*gin.Context"),
				tc.args.nickname, tc.args.email,
			).
				Return(&user.User{Nickname: tc.args.nickname, Email: tc.args.email}, tc.args.err).
				Once()
			var (
				response userResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.createUser(nil, nil)
			} else {
				response, err = suite.Client.createUser(tc.args.nickname, tc.args.email)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.nickname, response.Data.Nickname)
				suite.Equal(tc.args.email, response.Data.Email)
				suite.Equal(int64(0), response.Data.ID)
			}
		})
	}
}

func (suite *HTTPSuite) TestHandler_GetUser() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful get",
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:  1,
				err: app.ErrUserNotFound,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:  1,
				err: ErrMock,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.App.On("GetUser",
				mock.AnythingOfType("*gin.Context"),
				tc.args.id,
			).
				Return(&user.User{ID: tc.args.id, Nickname: "Mac Miller", Email: "swimming@circles.com"}, tc.args.err).
				Once()
			var (
				response userResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.getUser("hi")
			} else {
				response, err = suite.Client.getUser(tc.args.id)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal("Mac Miller", response.Data.Nickname)
				suite.Equal("swimming@circles.com", response.Data.Email)
				suite.Equal(tc.args.id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPSuite) TestHandler_UpdateUser() {
	type args struct {
		badId    bool
		badBody  bool
		id       int64
		nickname string
		email    string
		err      error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful update",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			wantErr: false,
		},
		{
			name: "validation error",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      validator.ValidationErrors{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      app.ErrUserNotFound,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      ErrMock,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
		{
			name: "bad request: id not int",
			args: args{
				badId:    true,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: unable to bind data",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				badBody:  true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "invalid email",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming.com",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			id, name, email, e := tc.args.id, tc.args.nickname, tc.args.email, tc.args.err
			suite.App.On("UpdateUser",
				mock.AnythingOfType("*gin.Context"),
				id, name, email,
			).
				Return(&user.User{ID: id, Nickname: name, Email: email}, e).
				Once()
			var (
				response userResponse
				err      error
			)
			if tc.args.badId {
				response, err = suite.Client.updateUser("hi", name, email)
			} else if tc.args.badBody {
				response, err = suite.Client.updateUser(id, 13, email)
			} else {
				response, err = suite.Client.updateUser(id, name, email)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(name, response.Data.Nickname)
				suite.Equal(email, response.Data.Email)
				suite.Equal(id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPSuite) TestHandler_DeleteUser() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful get",
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:  1,
				err: app.ErrUserNotFound,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:  1,
				err: ErrMock,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.App.On("DeleteUser",
				mock.AnythingOfType("*gin.Context"),
				tc.args.id,
			).
				Return(tc.args.err).
				Once()
			var err error
			if tc.args.badReq {
				_, err = suite.Client.deleteUser("hi")
			} else {
				_, err = suite.Client.deleteUser(tc.args.id)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPSuite))
}
