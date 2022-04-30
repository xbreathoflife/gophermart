package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"log"
	"os"
)

type DBStorage struct {
	ConnString string
}

func NewDBStorage(connString string) *DBStorage {
	storage := &DBStorage{ConnString: connString}
	return storage
}

func (s *DBStorage) connect(_ context.Context) (*sql.DB, error) {
	if s.ConnString == "" {
		log.Fatal("Connection string is empty\n")
	}
	conn, err := sql.Open("pgx", s.ConnString)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func (s *DBStorage) Init(ctx context.Context) error {
	// run migrations
	createTableQuery, err := os.ReadFile("./migrations/2022-04-16-create-tables.sql")
	if err != nil {
		return err
	}
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, string(createTableQuery))

	return err
}

// ------- INSERT QUERIES

func (s *DBStorage) InsertNewUser(ctx context.Context, user entities.UserModel) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO users(login, password_hash, session) VALUES ($1, $2, $3)`,
		user.Login, user.PasswordHash, user.Session)

	return err
}

func (s *DBStorage) UpdateUserSession(ctx context.Context, userSession entities.UserSessionModel) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE users SET session = $1 WHERE login = $2`,
		userSession.Session, userSession.Login)

	return err
}

func (s *DBStorage) InsertNewOrder(ctx context.Context, order entities.OrderModel) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO orders(order_num, login, uploaded_at, status) VALUES ($1, $2, $3, $4)`,
		order.OrderNum, order.Login, order.UploadedAt, order.Status)

	return err
}

func (s *DBStorage) InsertNewBalance(ctx context.Context, balance entities.BalanceModel) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO balance(login) VALUES ($1)`, balance.Login)

	return err
}

func (s *DBStorage) InsertNewBalanceWithdrawals(ctx context.Context, balanceWithdrawals entities.BalanceWithdrawalsModel) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO balance_withdrawals(login, order_num, sum, processed_at) VALUES ($1, $2, $3, $4)`,
		balanceWithdrawals.Login, balanceWithdrawals.OrderNum, balanceWithdrawals.Sum, balanceWithdrawals.ProcessedAt)

	return err
}

// ------- UPDATE QUERIES

func (s *DBStorage) UpdateBalance(ctx context.Context, balance entities.BalanceModel) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE balance SET balance = $1, spent = $2 WHERE login = $3`,
		balance.Balance, balance.Spent, balance.Login)

	return err
}

func (s *DBStorage) UpdateOrderStatus(ctx context.Context, orderNum string, status string) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE orders SET status = $1 WHERE order_num = $2`, status, orderNum)

	return err
}

func (s *DBStorage) UpdateOrderStatusAndAccrual(ctx context.Context, orderNum string, status string, accrual float64) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE order_num = $3`,
		status, accrual, orderNum)

	return err
}

// ------- GET QUERIES

func (s *DBStorage) GetUserIfExists(ctx context.Context, login string) (*entities.UserModel, error) {
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	var user entities.UserModel
	row := conn.QueryRowContext(ctx,
		`SELECT login, password_hash FROM users WHERE login = $1`, login)
	err = row.Scan(&user.Login, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *DBStorage) GetUserBySessionIfExists(ctx context.Context, session string) (*entities.UserSessionModel, error) {
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	var userSession entities.UserSessionModel
	row := conn.QueryRowContext(ctx,
		`SELECT login, session FROM users WHERE session = $1`, session)
	err = row.Scan(&userSession.Login, &userSession.Session)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &userSession, nil
}

func (s *DBStorage) GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderModel, error) {
	conn, err := s.connect(ctx)
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

func (s *DBStorage) GetOrderIfExists(ctx context.Context, orderNum string) (*entities.OrderModel, error) {
	conn, err := s.connect(ctx)
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

func (s *DBStorage) GetBalance(ctx context.Context, login string) (*entities.BalanceModel, error) {
	conn, err := s.connect(ctx)
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

func (s *DBStorage) GetBalanceWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsModel, error) {
	conn, err := s.connect(ctx)
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
