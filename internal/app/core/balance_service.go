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
)

type BalanceService struct {
	Storage storage.Storage
}

func NewBalanceService(storage storage.Storage) *BalanceService {
	service := BalanceService{Storage: storage}
	return &service
}

func (bs *BalanceService) GetUsersBalance(ctx context.Context, login string) (*entities.BalanceModel, error) {
	balance, err := bs.Storage.GetBalance(ctx, login)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (bs *BalanceService) ProcessBalanceWithdraw(ctx context.Context, login string, bw entities.BalanceWithdrawRequest) (int, error) {
	if !luhn.Valid(bw.Order) {
		return http.StatusUnprocessableEntity, errors.NewWrongDataError(bw.Order)
	}

	balance, err := bs.Storage.GetBalance(ctx, login)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if balance.Balance < bw.Sum {
		return http.StatusPaymentRequired, errors2.New("not enough money")
	}

	err = bs.Storage.UpdateBalance(ctx, entities.BalanceModel{
		Login:   login,
		Balance: balance.Balance - bw.Sum,
		Spent:   balance.Spent + bw.Sum,
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = bs.Storage.InsertNewBalanceWithdrawals(ctx, entities.BalanceWithdrawalsModel{
		Login:       login,
		OrderNum:    bw.Order,
		Sum:         bw.Sum,
		ProcessedAt: time.Now(),
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (bs *BalanceService) GetWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsResponse, error) {
	withdrawalHistory, err := bs.Storage.GetBalanceWithdrawalsForUser(ctx, login)
	if err != nil {
		return nil, err
	}
	var bwResponse []entities.BalanceWithdrawalsResponse
	for _, bw := range withdrawalHistory {
		bwResponse = append(bwResponse, entities.BalanceWithdrawalsResponse{
			OrderNum:    bw.OrderNum,
			Sum:         bw.Sum,
			ProcessedAt: bw.ProcessedAt.Format(time.RFC3339),
		})
	}

	return bwResponse, nil
}
