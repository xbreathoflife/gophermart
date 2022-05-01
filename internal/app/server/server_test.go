package server

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"github.com/xbreathoflife/gophermart/internal/app/storage/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_RegisterLogin(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name     string
		userData entities.LoginRequest
		target   string
		want     want
	}{
		{
			name: "simple register",
			userData: entities.LoginRequest{Login: "hello", Password: "there"},
			target: "/api/user/register",
			want: want{
				statusCode: 200,
			},
		},
		{
			name: "register bad request",
			userData: entities.LoginRequest{Login: "hello"},
			target: "/api/user/register",
			want: want{
				statusCode: 400,
			},
		},
		{
			name: "register username taken",
			userData: entities.LoginRequest{Login: "goodbye", Password: "there"},
			target: "/api/user/register",
			want: want{
				statusCode: 409,
			},
		},
		{
			name: "simple login",
			userData: entities.LoginRequest{Login: "goodbye", Password: "123456"},
			target: "/api/user/login",
			want: want{
				statusCode: 200,
			},
		},
		{
			name: "login bad request",
			userData: entities.LoginRequest{Login: "goodbye"},
			target: "/api/user/login",
			want: want{
				statusCode: 400,
			},
		},
		{
			name: "login unauthorized",
			userData: entities.LoginRequest{Login: "goodbye", Password: "34"},
			target: "/api/user/login",
			want: want{
				statusCode: 401,
			},
		},

	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockStorage(mockCtrl)
	repo.EXPECT().Init(gomock.Any())
	repo.EXPECT().GetUserIfExists(gomock.Any(), gomock.Eq("hello")).Return(nil, nil).MinTimes(0)
	repo.EXPECT().GetUserIfExists(gomock.Any(), gomock.Eq("goodbye")).Return(
		&entities.UserModel{Login: "goodbye", PasswordHash: "123456", Session: "123"}, nil).MinTimes(0)
	repo.EXPECT().InsertNewUser(gomock.Any(), gomock.Any()).MinTimes(0)
	repo.EXPECT().InsertNewBalance(gomock.Any(), gomock.Any()).MinTimes(0)
	repo.EXPECT().UpdateUserSession(gomock.Any(), gomock.Any()).MinTimes(0)

	server := NewGothServer(repo, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.userData)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, tt.target, bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			h := server.ServerHandler()
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

		})
	}
}