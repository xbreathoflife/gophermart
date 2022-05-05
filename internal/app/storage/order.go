package storage

import (
	"context"
	"database/sql"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
)

type OrderStorage interface {
	InsertNewOrder(ctx context.Context, order entities.OrderModel) error
	UpdateOrderStatus(ctx context.Context, orderNum string, status string) error
	UpdateOrderStatusAndAccrual(ctx context.Context, orderNum string, status string, accrual float64) error
	GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderModel, error)
	GetOrderIfExists(ctx context.Context, orderNum string) (*entities.OrderModel, error)
}

type OrderStorageImpl struct {
	ConnString string
}

func NewOrderStorage(connString string) *OrderStorageImpl {
	storage := &OrderStorageImpl{ConnString: connString}
	return storage
}

func (s *OrderStorageImpl) InsertNewOrder(ctx context.Context, order entities.OrderModel) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO orders(order_num, login, uploaded_at, status) VALUES ($1, $2, $3, $4)`,
		order.OrderNum, order.Login, order.UploadedAt, order.Status)

	return err
}

func (s *OrderStorageImpl) UpdateOrderStatus(ctx context.Context, orderNum string, status string) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE orders SET status = $1 WHERE order_num = $2`, status, orderNum)

	return err
}

func (s *OrderStorageImpl) UpdateOrderStatusAndAccrual(ctx context.Context, orderNum string, status string, accrual float64) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE order_num = $3`,
		status, accrual, orderNum)

	return err
}

func (s *OrderStorageImpl) GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderModel, error) {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	rows, err := conn.QueryContext(ctx,
		`SELECT order_num, login, uploaded_at, status, accrual FROM orders
				WHERE login = $1 ORDER BY uploaded_at`, login)
	if err != nil && rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []entities.OrderModel
	for rows.Next() {
		var o entities.OrderModel
		if err := rows.Scan(&o.OrderNum, &o.Login, &o.UploadedAt, &o.Status, &o.Accrual); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (s *OrderStorageImpl) GetOrderIfExists(ctx context.Context, orderNum string) (*entities.OrderModel, error) {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	var order entities.OrderModel
	row := conn.QueryRowContext(ctx,
		`SELECT order_num, login, uploaded_at, status FROM orders
				WHERE order_num = $1`, orderNum)
	err = row.Scan(&order.OrderNum, &order.Login, &order.UploadedAt, &order.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}
