package storage

import (
	"context"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
)

type Storage interface {
	Init(ctx context.Context) error
	InsertNewUser(ctx context.Context, user entities.UserModel) error
	UpdateUserSession(ctx context.Context, userSession entities.UserSessionModel) error
	InsertNewOrder(ctx context.Context, order entities.OrderModel) error
	InsertNewBalance(ctx context.Context, balance entities.BalanceModel) error
	InsertNewBalanceWithdrawals(ctx context.Context, balanceWithdrawals entities.BalanceWithdrawalsModel) error
	UpdateBalance(ctx context.Context, balance entities.BalanceModel) error
	UpdateOrderStatus(ctx context.Context, orderNum string, status string) error
	UpdateOrderStatusAndAccrual(ctx context.Context, orderNum string, status string, accrual float64) error
	GetUserIfExists(ctx context.Context, login string) (*entities.UserModel, error)
	GetUserBySessionIfExists(ctx context.Context, session string) (*entities.UserSessionModel, error)
	GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderModel, error)
	GetOrderIfExists(ctx context.Context, orderNum string) (*entities.OrderModel, error)
	GetBalance(ctx context.Context, login string) (*entities.BalanceModel, error)
	GetBalanceWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsModel, error)
}
