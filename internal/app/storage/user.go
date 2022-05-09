package storage

import (
	"context"
	"database/sql"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
)

type UserStorage interface {
	InsertNewUser(ctx context.Context, user entities.UserModel) error
	UpdateUserSession(ctx context.Context, userSession entities.UserSessionModel) error
	GetUserIfExists(ctx context.Context, login string) (*entities.UserModel, error)
	GetUserBySessionIfExists(ctx context.Context, session string) (*entities.UserSessionModel, error)
}

type UserStorageImpl struct {
	ConnString string
}

func NewUserStorage(connString string) *UserStorageImpl {
	storage := &UserStorageImpl{ConnString: connString}
	return storage
}

func (s *UserStorageImpl) InsertNewUser(ctx context.Context, user entities.UserModel) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO users(login, password_hash, session) VALUES ($1, $2, $3)`,
		user.Login, user.PasswordHash, user.Session)

	return err
}

func (s *UserStorageImpl) UpdateUserSession(ctx context.Context, userSession entities.UserSessionModel) error {
	conn, err := connect(ctx, s.ConnString)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`UPDATE users SET session = $1 WHERE login = $2`,
		userSession.Session, userSession.Login)

	return err
}

func (s *UserStorageImpl) GetUserIfExists(ctx context.Context, login string) (*entities.UserModel, error) {
	conn, err := connect(ctx, s.ConnString)
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

func (s *UserStorageImpl) GetUserBySessionIfExists(ctx context.Context, session string) (*entities.UserSessionModel, error) {
	conn, err := connect(ctx, s.ConnString)
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