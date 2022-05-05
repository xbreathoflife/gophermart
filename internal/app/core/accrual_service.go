package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	RegisteredStatus = "REGISTERED"
	ProcessingStatus = "PROCESSING"
	InvalidStatus    = "INVALID"
	ProcessedStatus  = "PROCESSED"
)

type AccrualService struct {
	Storage        storage.Storage
	ServiceAddress string
	Channel        chan string
	isNotFinished  bool
}

func NewAccrualService(storage storage.Storage, serviceAddress string, ctx context.Context) *AccrualService {
	ch := make(chan string, 10)
	service := AccrualService{Storage: storage, ServiceAddress: serviceAddress, Channel: ch, isNotFinished: true}
	if serviceAddress != "" {
		go service.updateOrderStatuses(ctx)
	}
	return &service
}

func (as *AccrualService) updateOrderStatuses(ctx context.Context) {
	for as.isNotFinished {
		orderNum := <-as.Channel
		fmt.Printf("%s\n", orderNum)
		resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.ServiceAddress, orderNum))
		if err != nil {
			log.Printf("Failed to get %s/api/orders/%s\n", as.ServiceAddress, orderNum)
		}

		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			log.Println("Too many requests status")
			time.Sleep(time.Second * 3)
			as.Channel <- orderNum
		case http.StatusInternalServerError:
			log.Println("Something went wrong 500 status for order ", orderNum)
			as.isNotFinished = false
		case http.StatusOK:
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}

			orderStatus := entities.GetOrderStatusResponse{}
			if err := json.Unmarshal(b, &orderStatus); err != nil {
				log.Println("Failed to parse json", err)
				log.Println("body: ", b)
				return
			}
			switch orderStatus.Status {
			case InvalidStatus:
				log.Println("Invalid status for order ", orderNum)
				err := as.Storage.UpdateOrderStatus(ctx, orderNum, InvalidStatus)
				if err != nil {
					log.Println("Failed to update status in db: ", err)
				}
			case ProcessedStatus:
				log.Println("Processed status for order ", orderNum)
				accrual := *orderStatus.Accrual
				err := as.Storage.UpdateOrderStatusAndAccrual(ctx, orderNum, ProcessedStatus, accrual)
				if err != nil {
					log.Println("Failed to update status in db: ", err)
				}
				// get user login
				o, err := as.Storage.GetOrderIfExists(ctx, orderNum)
				login := o.Login

				if err != nil {
					log.Println("Failed to update status in db: ", err)
				}
				balance, err := as.Storage.GetBalance(ctx, login)
				if err != nil {
					log.Println("Failed to update status in db: ", err)
				}

				err = as.Storage.UpdateBalance(ctx, entities.BalanceModel{
					Login:   login,
					Balance: balance.Balance + accrual,
					Spent:   balance.Spent,
				})
				if err != nil {
					log.Println("Failed to update status in db: ", err)
				}
			case ProcessingStatus:
				log.Println("Processing status for order ", orderNum)
				err := as.Storage.UpdateOrderStatus(ctx, orderNum, ProcessingStatus)
				if err != nil {
					log.Println("Failed to update status in db: ", err)
				}
				as.Channel <- orderNum
			case RegisteredStatus:
				log.Println("New status for order ", orderNum)
				as.Channel <- orderNum
			}
		}
	}
}
