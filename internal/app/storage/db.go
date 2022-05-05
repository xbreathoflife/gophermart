package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"os"
)

type CommonStorage interface {
	Init(ctx context.Context) error
}

type DBStorage struct {
	ConnString string
}

func NewDBStorage(connString string) *DBStorage {
	storage := &DBStorage{ConnString: connString}
	return storage
}

func connect(_ context.Context, connString string) (*sql.DB, error) {
	if connString == "" {
		log.Fatal("Connection string is empty\n")
	}
	conn, err := sql.Open("pgx", connString)
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
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, string(createTableQuery))

	return err
}