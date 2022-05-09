package core

import (
	"context"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"github.com/xbreathoflife/gophermart/internal/app/errors"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
)

type UserService struct {
	UserStorage    storage.UserStorage
	BalanceStorage storage.BalanceStorage
}

func NewUserService(userStorage storage.UserStorage, balanceStorage storage.BalanceStorage) *UserService {
	service := UserService{UserStorage: userStorage, BalanceStorage: balanceStorage}
	return &service
}

func (us *UserService) CheckUserExists(ctx context.Context, user entities.LoginRequest) error {
	prevUser, err := us.UserStorage.GetUserIfExists(ctx, user.Login)
	if err != nil {
		return err
	}
	if prevUser != nil {
		return errors.NewDuplicateError(prevUser.Login)
	}
	return nil
}

func (us *UserService) InsertNewUser(ctx context.Context, user entities.UserModel) error {
	err := us.UserStorage.InsertNewUser(ctx, user)
	if err != nil {
		return err
	}
	return us.BalanceStorage.InsertNewBalance(ctx, entities.BalanceModel{Login: user.Login})
}

func (us *UserService) CheckUserCredentials(ctx context.Context, user entities.LoginRequest) error {
	prevUser, err := us.UserStorage.GetUserIfExists(ctx, user.Login)
	if err != nil {
		return err
	}
	if prevUser == nil || (prevUser.PasswordHash != user.Password || prevUser.Login != user.Login) {
		return errors.NewWrongDataError(prevUser.Login)
	}

	return nil
}

func (us *UserService) UpdateUserSession(ctx context.Context, userSession entities.UserSessionModel) error {
	return us.UserStorage.UpdateUserSession(ctx, userSession)
}

func (us *UserService) GetUserBySession(ctx context.Context, session string) (*entities.UserSessionModel, error) {
	sessionModel, err := us.UserStorage.GetUserBySessionIfExists(ctx, session)
	if err != nil {
		return nil, err
	}
	if sessionModel == nil {
		return nil, errors.NewWrongDataError(session)
	}
	return sessionModel, nil
}
