package core

import (
	"context"
	errors2 "errors"
	"github.com/joeljunstrom/go-luhn"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"github.com/xbreathoflife/gophermart/internal/app/errors"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
	"net/http"
	"time"
	"unicode"
)

const NewStatus = "NEW"

type OrderService struct {
	OrderStorage storage.OrderStorage
	Accrual      *AccrualService
}

func NewOrderService(orderStorage storage.OrderStorage, balanceStorage storage.BalanceStorage, serviceAddress string, ctx context.Context) *OrderService {
	accrual := NewAccrualService(orderStorage, balanceStorage, serviceAddress, ctx)
	service := OrderService{OrderStorage: orderStorage, Accrual: accrual}
	return &service
}

func IsNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func (os *OrderService) CreateNewOrder(ctx context.Context, login string, orderNum string) (int, error) {
	if !IsNumber(orderNum) {
		return http.StatusBadRequest, errors2.New("not a number")
	}

	if !luhn.Valid(orderNum) {
		return http.StatusUnprocessableEntity, errors.NewWrongDataError(orderNum)
	}

	order, err := os.OrderStorage.GetOrderIfExists(ctx, orderNum)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if order != nil {
		if order.Login == login {
			return http.StatusOK, nil
		}
		return http.StatusConflict, errors.NewDuplicateError(orderNum)
	}

	err = os.OrderStorage.InsertNewOrder(ctx, entities.OrderModel{
		OrderNum:   orderNum,
		Login:      login,
		UploadedAt: time.Now(),
		Status:     NewStatus,
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	os.Accrual.Channel <- orderNum // отправляем результат в канал

	return http.StatusAccepted, nil
}

func (os *OrderService) GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderResponse, error) {
	orders, err := os.OrderStorage.GetOrdersForUser(ctx, login)
	if err != nil {
		return nil, err
	}
	var ordersResponse []entities.OrderResponse
	for _, o := range orders {
		var accrual *float64 = nil
		if o.Accrual.Valid {
			accrual = &(o.Accrual.Float64)
		}
		ordersResponse = append(ordersResponse, entities.OrderResponse{
			OrderNum:   o.OrderNum,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
			Status:     o.Status,
			Accrual:    accrual,
		})
	}

	return ordersResponse, nil
}
