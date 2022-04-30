package core

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/joeljunstrom/go-luhn"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"github.com/xbreathoflife/gophermart/internal/app/errors"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	NewStatus        = "NEW"
	RegisteredStatus = "REGISTERED"
	ProcessingStatus = "PROCESSING"
	InvalidStatus    = "INVALID"
	ProcessedStatus  = "PROCESSED"
)

type LoyaltyService struct {
	Storage        storage.DBStorage
	Channel        chan string
	ServiceAddress string
}

func NewLoyaltyService(storage storage.DBStorage, serviceAddress string, ctx context.Context) *LoyaltyService {
	ch := make(chan string, 10)
	service := LoyaltyService{Storage: storage, Channel: ch, ServiceAddress: serviceAddress}
	go service.updateOrderStatuses(ctx)
	return &service
}

func (ls *LoyaltyService) updateOrderStatuses(ctx context.Context) {
	for {
		orderNum := <-ls.Channel
		fmt.Printf("%s\n", orderNum)
		resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", ls.ServiceAddress, orderNum))
		if err != nil {
			log.Fatalln(err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("Too many requests status")
			time.Sleep(time.Second * 3)
			ls.Channel <- orderNum
		} else if resp.StatusCode == http.StatusInternalServerError {
			log.Fatalln("Something went wrong 500 status for order", orderNum)
			// todo: finish??? invalid????
		} else if resp.StatusCode == http.StatusOK {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			orderStatus := entities.GetOrderStatusResponse{}
			if err := json.Unmarshal(b, &orderStatus); err != nil {
				log.Fatalln("Failed to parse json")
				return
			}
			if orderStatus.Status == InvalidStatus {
				log.Println("Invalid status for order ", orderNum)
				err := ls.Storage.UpdateOrderStatus(ctx, orderNum, InvalidStatus)
				if err != nil {
					log.Fatalln("Failed to update status in db: ", err)
				}
			} else if orderStatus.Status == ProcessedStatus {
				log.Println("Processed status for order ", orderNum)
				accrual := *orderStatus.Accrual
				err := ls.Storage.UpdateOrderStatusAndAccrual(ctx, orderNum, ProcessedStatus, accrual)
				if err != nil {
					log.Fatalln("Failed to update status in db: ", err)
				}
				// get user login
				o, err := ls.Storage.GetOrderIfExists(ctx, orderNum)
				login := o.Login

				if err != nil {
					log.Fatalln("Failed to update status in db: ", err)
				}
				balance, err := ls.Storage.GetBalance(ctx, login)
				if err != nil {
					log.Fatalln("Failed to update status in db: ", err)
				}

				err = ls.Storage.UpdateBalance(ctx, entities.BalanceModel{
					Login:   login,
					Balance: balance.Balance + accrual,
					Spent:   balance.Spent,
				})
				if err != nil {
					log.Fatalln("Failed to update status in db: ", err)
				}
			} else if orderStatus.Status == ProcessingStatus {
				log.Println("Processing status for order ", orderNum)
				err := ls.Storage.UpdateOrderStatus(ctx, orderNum, ProcessingStatus)
				if err != nil {
					log.Fatalln("Failed to update status in db: ", err)
				}
				time.Sleep(time.Second)
				ls.Channel <- orderNum
			} else if orderStatus.Status == RegisteredStatus {
				log.Println("New status for order ", orderNum)
				time.Sleep(time.Second)
				ls.Channel <- orderNum
			}
		}
	}
}

func (ls *LoyaltyService) CheckUserExists(ctx context.Context, user entities.LoginRequest) error {
	prevUser, err := ls.Storage.GetUserIfExists(ctx, user.Login)
	if err != nil {
		return err
	}
	if prevUser != nil {
		return errors.NewDuplicateError(prevUser.Login)
	}
	return nil
}

func (ls *LoyaltyService) InsertNewUser(ctx context.Context, user entities.UserModel) error {
	err := ls.Storage.InsertNewUser(ctx, user)
	if err != nil {
		return err
	}
	return ls.Storage.InsertNewBalance(ctx, entities.BalanceModel{Login: user.Login})
}

func (ls *LoyaltyService) CheckUserCredentials(ctx context.Context, user entities.LoginRequest) error {
	prevUser, err := ls.Storage.GetUserIfExists(ctx, user.Login)
	if err != nil {
		return err
	}
	if prevUser != nil || (prevUser.PasswordHash != user.Password && prevUser.Login != user.Login) {
		return errors.NewWrongDataError(prevUser.Login)
	}

	return nil
}

func (ls *LoyaltyService) UpdateUserSession(ctx context.Context, userSession entities.UserSessionModel) error {
	return ls.Storage.UpdateUserSession(ctx, userSession)
}

func (ls *LoyaltyService) GetUserBySession(ctx context.Context, session string) (*entities.UserSessionModel, error) {
	sessionModel, err := ls.Storage.GetUserBySessionIfExists(ctx, session)
	if err != nil {
		return nil, err
	}
	if sessionModel == nil {
		return nil, errors.NewWrongDataError(session)
	}
	return sessionModel, nil
}

func (ls *LoyaltyService) CreateNewOrder(ctx context.Context, login string, orderNum string) (error, int) {
	if !luhn.Valid(orderNum) {
		return errors.NewWrongDataError(orderNum), http.StatusUnprocessableEntity
	}

	order, err := ls.Storage.GetOrderIfExists(ctx, orderNum)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if order != nil {
		if order.Login == login {
			return nil, http.StatusOK
		}
		return errors.NewDuplicateError(orderNum), http.StatusConflict
	}

	err = ls.Storage.InsertNewOrder(ctx, entities.OrderModel{
		OrderNum:   orderNum,
		Login:      login,
		UploadedAt: time.Now(),
		Status:     NewStatus,
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}

	// todo: add to queue
	ls.Channel <- orderNum // отправляем результат в канал

	return nil, http.StatusAccepted
}

func (ls *LoyaltyService) GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderResponse, error) {
	orders, err := ls.Storage.GetOrdersForUser(ctx, login)
	if err != nil {
		return nil, err
	}
	var ordersResponse []entities.OrderResponse
	for _, o := range orders {
		var accrual *int64 = nil
		if o.Accrual.Valid {
			accrual = &o.Accrual.Int64
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

func (ls *LoyaltyService) GetUsersBalance(ctx context.Context, login string) (*entities.BalanceModel, error) {
	balance, err := ls.Storage.GetBalance(ctx, login)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (ls *LoyaltyService) ProcessBalanceWithdraw(ctx context.Context, login string, bw entities.BalanceWithdrawRequest) (error, int) {
	if !luhn.Valid(bw.Order) {
		return errors.NewWrongDataError(bw.Order), http.StatusUnprocessableEntity
	}

	balance, err := ls.Storage.GetBalance(ctx, login)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	if balance.Balance < bw.Sum {
		return errors2.New("Not enough money."), http.StatusPaymentRequired
	}

	err = ls.Storage.UpdateBalance(ctx, entities.BalanceModel{
		Login:   login,
		Balance: balance.Balance - bw.Sum,
		Spent:   balance.Spent + bw.Sum,
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = ls.Storage.InsertNewBalanceWithdrawals(ctx, entities.BalanceWithdrawalsModel{
		Login:       login,
		OrderNum:    bw.Order,
		Sum:         bw.Sum,
		ProcessedAt: time.Now(),
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (ls *LoyaltyService) GetWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsResponse, error) {
	withdrawalHistory, err := ls.Storage.GetBalanceWithdrawalsForUser(ctx, login)
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
