package entities

import "time"

type UserModel struct {
	Login        string
	PasswordHash string
}

type UserSessionModel struct {
	Session string
	Login   string
}

type OrderModel struct {
	OrderNum   int
	Login      string
	UploadedAt time.Time
	Status     string
	Accrual    int
}

type BalanceModel struct {
	Login   string
	Balance int
	Spent   int
}

type BalanceWithdrawalsModel struct {
	Login       string
	OrderNum    int
	Sum         int
	ProcessedAt time.Time
}