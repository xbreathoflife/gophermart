package entities

import (
	"database/sql"
	"time"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type OrderResponse struct {
	OrderNum   string   `json:"number"`
	UploadedAt string   `json:"uploaded_at"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
}

type BalanceWithdrawalsResponse struct {
	OrderNum    string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type UserModel struct {
	Login        string
	PasswordHash string
	Session      string
}

type UserSessionModel struct {
	Session string
	Login   string
}

type OrderModel struct {
	OrderNum   string
	Login      string
	UploadedAt time.Time
	Status     string
	Accrual    sql.NullFloat64
}

type BalanceModel struct {
	Login   string  `json:"-"`
	Balance float64 `json:"current"`
	Spent   float64 `json:"withdrawn"`
}

type BalanceWithdrawalsModel struct {
	Login       string
	OrderNum    string
	Sum         float64
	ProcessedAt time.Time
}

type GetOrderStatusResponse struct {
	OrderNum string   `json:"order"`
	Status   string   `json:"status"`
	Accrual  *float64 `json:"accrual,omitempty"`
}
