package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"github.com/xbreathoflife/gophermart/internal/app/storage/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
			name:     "simple register",
			userData: entities.LoginRequest{Login: "hello", Password: "there"},
			target:   "/api/user/register",
			want: want{
				statusCode: 200,
			},
		},
		{
			name:     "register bad request",
			userData: entities.LoginRequest{Login: "hello"},
			target:   "/api/user/register",
			want: want{
				statusCode: 400,
			},
		},
		{
			name:     "register username taken",
			userData: entities.LoginRequest{Login: "goodbye", Password: "there"},
			target:   "/api/user/register",
			want: want{
				statusCode: 409,
			},
		},
		{
			name:     "simple login",
			userData: entities.LoginRequest{Login: "goodbye", Password: "123456"},
			target:   "/api/user/login",
			want: want{
				statusCode: 200,
			},
		},
		{
			name:     "login bad request",
			userData: entities.LoginRequest{Login: "goodbye"},
			target:   "/api/user/login",
			want: want{
				statusCode: 400,
			},
		},
		{
			name:     "login unauthorized",
			userData: entities.LoginRequest{Login: "goodbye", Password: "34"},
			target:   "/api/user/login",
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

			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func checkAuth(server *gophServer, t *testing.T) *http.Cookie {
	body, err := json.Marshal(entities.LoginRequest{Login: "hello", Password: "123456"})
	require.NoError(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h := server.ServerHandler()
	h.ServeHTTP(w, request)
	result := w.Result()
	auth := result.Cookies()[0]
	cookie := http.Cookie{Name: auth.Name, Value: auth.Value}

	err = result.Body.Close()
	require.NoError(t, err)
	return &cookie
}

func TestServer_InsertNewOrder(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name     string
		orderNum string
		target   string
		want     want
	}{
		{
			name:     "add new order",
			orderNum: "2377225624",
			target:   "/api/user/orders",
			want: want{
				statusCode: 202,
			},
		},
		{
			name:     "failed luhn check",
			orderNum: "123",
			target:   "/api/user/orders",
			want: want{
				statusCode: 422,
			},
		},
		{
			name:     "bad request",
			orderNum: "heh",
			target:   "/api/user/orders",
			want: want{
				statusCode: 400,
			},
		},
		{
			name:     "bad request",
			orderNum: "562246784655",
			target:   "/api/user/orders",
			want: want{
				statusCode: 409,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockStorage(mockCtrl)
	repo.EXPECT().Init(gomock.Any())
	repo.EXPECT().GetUserIfExists(gomock.Any(), gomock.Eq("hello")).Return(
		&entities.UserModel{Login: "hello", PasswordHash: "123456", Session: "123"}, nil).MinTimes(0)
	repo.EXPECT().UpdateUserSession(gomock.Any(), gomock.Any()).MinTimes(0)
	repo.EXPECT().GetUserBySessionIfExists(gomock.Any(), gomock.Any()).Return(
		&entities.UserSessionModel{Login: "hello", Session: "123"}, nil).MinTimes(0)
	repo.EXPECT().GetOrderIfExists(gomock.Any(), "2377225624").Return(nil, nil).MinTimes(0)

	orderTime, err := time.Parse(time.RFC3339, "2022-04-30T20:00:00+03:00")
	require.NoError(t, err)

	repo.EXPECT().GetOrderIfExists(gomock.Any(), "12345678903").Return(
		&entities.OrderModel{OrderNum: "12345678903",
			Login:      "hello",
			UploadedAt: orderTime,
			Status:     "NEW",
		}, nil).MinTimes(0)
	repo.EXPECT().GetOrderIfExists(gomock.Any(), "562246784655").Return(
		&entities.OrderModel{OrderNum: "562246784655",
			Login:      "goodbye",
			UploadedAt: orderTime,
			Status:     "NEW",
		}, nil).MinTimes(0)
	repo.EXPECT().InsertNewOrder(gomock.Any(), gomock.Any()).MinTimes(0)

	server := NewGothServer(repo, "")
	cookie := checkAuth(server, t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(tt.orderNum)
			request := httptest.NewRequest(http.MethodPost, tt.target, bytes.NewBuffer(body))
			request.AddCookie(cookie)
			w := httptest.NewRecorder()
			h := server.ServerHandler()
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestServer_GetOrders(t *testing.T) {
	orderTime, err := time.Parse(time.RFC3339, "2022-04-30T20:00:00+03:00")
	require.NoError(t, err)
	orders := []entities.OrderModel{
		{
			OrderNum:   "12345674",
			Login:      "hello",
			UploadedAt: orderTime,
			Status:     "NEW",
		},
		{
			OrderNum:   "12345678903",
			Login:      "hello",
			UploadedAt: orderTime,
			Status:     "PROCESSING",
		},
		{
			OrderNum:   "9278923470",
			Login:      "hello",
			UploadedAt: orderTime,
			Status:     "PROCESSED",
			Accrual:    sql.NullFloat64{Float64: 500, Valid: true},
		},
		{
			OrderNum:   "346436439",
			Login:      "hello",
			UploadedAt: orderTime,
			Status:     "INVALID",
			Accrual:    sql.NullFloat64{Float64: 500, Valid: true},
		},
	}
	var expected []entities.OrderResponse
	for _, o := range orders {
		var accrual *float64 = nil
		if o.Accrual.Valid {
			accrual = &(o.Accrual.Float64)
		}
		expected = append(expected, entities.OrderResponse{
			OrderNum:   o.OrderNum,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
			Status:     o.Status,
			Accrual:    accrual,
		})
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockStorage(mockCtrl)
	repo.EXPECT().Init(gomock.Any())
	repo.EXPECT().GetUserIfExists(gomock.Any(), gomock.Eq("hello")).Return(
		&entities.UserModel{Login: "hello", PasswordHash: "123456", Session: "123"}, nil).MinTimes(0)
	repo.EXPECT().UpdateUserSession(gomock.Any(), gomock.Any()).MinTimes(0)
	repo.EXPECT().GetUserBySessionIfExists(gomock.Any(), gomock.Any()).Return(
		&entities.UserSessionModel{Login: "hello", Session: "123"}, nil).MinTimes(0)

	repo.EXPECT().GetOrdersForUser(gomock.Any(), "hello").Return(
		orders, nil).MinTimes(0)

	server := NewGothServer(repo, "")
	cookie := checkAuth(server, t)

	request := httptest.NewRequest(http.MethodGet, "/api/user/orders", bytes.NewBuffer(nil))
	request.AddCookie(cookie)
	w := httptest.NewRecorder()
	h := server.ServerHandler()
	h.ServeHTTP(w, request)
	result := w.Result()

	assert.Equal(t, 200, result.StatusCode)

	ordersResult, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err)
	var actualOrders []entities.OrderResponse
	err = json.Unmarshal(ordersResult, &actualOrders)
	require.NoError(t, err)
	assert.ElementsMatch(t, actualOrders, expected)

	err = result.Body.Close()
	require.NoError(t, err)
}

func TestServer_GetBalance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockStorage(mockCtrl)
	repo.EXPECT().Init(gomock.Any())
	repo.EXPECT().GetUserIfExists(gomock.Any(), gomock.Eq("hello")).Return(
		&entities.UserModel{Login: "hello", PasswordHash: "123456", Session: "123"}, nil).MinTimes(0)
	repo.EXPECT().UpdateUserSession(gomock.Any(), gomock.Any()).MinTimes(0)
	repo.EXPECT().GetUserBySessionIfExists(gomock.Any(), gomock.Any()).Return(
		&entities.UserSessionModel{Login: "hello", Session: "123"}, nil).MinTimes(0)

	expectedBalance := entities.BalanceModel{Login: "hello", Balance: 510.5, Spent: 330}
	repo.EXPECT().GetBalance(gomock.Any(), gomock.Eq("hello")).Return(&expectedBalance, nil)

	server := NewGothServer(repo, "")
	cookie := checkAuth(server, t)

	request := httptest.NewRequest(http.MethodGet, "/api/user/balance", bytes.NewBuffer(nil))
	request.AddCookie(cookie)
	w := httptest.NewRecorder()
	h := server.ServerHandler()
	h.ServeHTTP(w, request)
	result := w.Result()

	assert.Equal(t, 200, result.StatusCode)

	balanceResult, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err)
	var actualBalance entities.BalanceModel
	err = json.Unmarshal(balanceResult, &actualBalance)
	require.NoError(t, err)
	expectedBalance.Login = ""
	assert.Equal(t, actualBalance, expectedBalance)

	err = result.Body.Close()
	require.NoError(t, err)
}


func TestServer_Withdrawals(t *testing.T) {
	processedAtStr := "2022-04-30T20:00:00+03:00"
	processedAt, err := time.Parse(time.RFC3339, processedAtStr)
	require.NoError(t, err)
	type want struct {
		statusCode  int
		withdrawals []entities.BalanceWithdrawalsResponse
	}
	tests := []struct {
		name    string
		request entities.BalanceWithdrawRequest
		target  string
		want    want
	}{
		{
			name:    "withdraw balance",
			request: entities.BalanceWithdrawRequest{Order: "2377225624", Sum: 500.5},
			target:  "/api/user/balance/withdraw",
			want: want{
				statusCode: 200,
				withdrawals: []entities.BalanceWithdrawalsResponse{
					{OrderNum: "2377225624", Sum: 500.5, ProcessedAt: processedAtStr},
				},
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockStorage(mockCtrl)
	repo.EXPECT().Init(gomock.Any())
	repo.EXPECT().GetUserIfExists(gomock.Any(), gomock.Eq("hello")).Return(
		&entities.UserModel{Login: "hello", PasswordHash: "123456", Session: "123"}, nil).MinTimes(0)
	repo.EXPECT().UpdateUserSession(gomock.Any(), gomock.Any()).MinTimes(0)
	repo.EXPECT().GetUserBySessionIfExists(gomock.Any(), gomock.Any()).Return(
		&entities.UserSessionModel{Login: "hello", Session: "123"}, nil).MinTimes(0)

	expectedBalance := entities.BalanceModel{Login: "hello", Balance: 510.5, Spent: 330}
	repo.EXPECT().GetBalance(gomock.Any(), gomock.Eq("hello")).Return(&expectedBalance, nil)
	repo.EXPECT().UpdateBalance(gomock.Any(), entities.BalanceModel{Login: "hello", Balance: 10, Spent: 830.5})
	repo.EXPECT().InsertNewBalanceWithdrawals(gomock.Any(), gomock.Any())
	repo.EXPECT().GetBalanceWithdrawalsForUser(gomock.Any(), "hello").Return(
		[]entities.BalanceWithdrawalsModel{{Login: "hello", OrderNum: "2377225624", Sum: 500.5, ProcessedAt: processedAt}}, nil)

	server := NewGothServer(repo, "")
	cookie := checkAuth(server, t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, tt.target, bytes.NewBuffer(body))
			request.AddCookie(cookie)
			w := httptest.NewRecorder()
			h := server.ServerHandler()
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err = result.Body.Close()
			require.NoError(t, err)
			request = httptest.NewRequest(http.MethodGet, "/api/user/balance/withdrawals", bytes.NewBuffer(nil))
			request.AddCookie(cookie)
			w = httptest.NewRecorder()
			h = server.ServerHandler()
			h.ServeHTTP(w, request)
			result = w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			wdResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			var actualWd []entities.BalanceWithdrawalsResponse
			err = json.Unmarshal(wdResult, &actualWd)
			require.NoError(t, err)
			assert.Equal(t, actualWd, tt.want.withdrawals)
		})
	}
}