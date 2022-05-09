package storage

import (
	"context"
	"database/sql"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
)

type BalanceStorage interface {
	InsertNewBalance(ctx context.Context, balance entities.BalanceModel) error
	InsertNewBalanceWithdrawals(ctx context.Context, balanceWithdrawals entities.BalanceWithdrawalsModel) error
	UpdateBalance(ctx context.Context, balance entities.BalanceModel) error
	GetBalance(ctx context.Context, login string) (*entities.BalanceModel, error)
	GetBalanceWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsModel, error)
}

type BalanceStorageImpl struct {
	ConnString string
}

func NewBalanceStorage(connString string) *BalanceStorageImpl {
	storage := &BalanceStorageImpl{ConnString: connString}
	return storage
}

func (s *BalanceStorageImpl) InsertNewBalance(ctx context.Context, balance entities.BalanceModel) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO balance(login) VALUES ($1)`, balance.Login)

	return err
}

func (s *BalanceStorageImpl) InsertNewBalanceWithdrawals(ctx context.Context, balanceWithdrawals entities.BalanceWithdrawalsModel) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO balance_withdrawals(login, order_num, sum, processed_at) VALUES ($1, $2, $3, $4)`,
		balanceWithdrawals.Login, balanceWithdrawals.OrderNum, balanceWithdrawals.Sum, balanceWithdrawals.ProcessedAt)

	return err
}

func (s *BalanceStorageImpl) UpdateBalance(ctx context.Context, balance entities.BalanceModel) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE balance SET balance = $1, spent = $2 WHERE login = $3`,
		balance.Balance, balance.Spent, balance.Login)

	return err
}

func (s *BalanceStorageImpl) GetBalance(ctx context.Context, login string) (*entities.BalanceModel, error) {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	var balance entities.BalanceModel
	row := conn.QueryRowContext(ctx,
		`SELECT login, balance, spent FROM balance WHERE login = $1`, login)
	err = row.Scan(&balance.Login, &balance.Balance, &balance.Spent)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &balance, nil
}

func (s *BalanceStorageImpl) GetBalanceWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsModel, error) {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	rows, err := conn.QueryContext(ctx,
		`SELECT login, order_num, sum, processed_at FROM balance_withdrawals 
				WHERE login = $1 ORDER BY processed_at`, login)
	if err != nil && rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	var balances []entities.BalanceWithdrawalsModel
	for rows.Next() {
		var b entities.BalanceWithdrawalsModel
		if err := rows.Scan(&b.Login, &b.OrderNum, &b.Sum, &b.ProcessedAt); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}

	return balances, nil
}
