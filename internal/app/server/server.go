package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xbreathoflife/gophermart/internal/app/auth"
	"github.com/xbreathoflife/gophermart/internal/app/core"
	"github.com/xbreathoflife/gophermart/internal/app/handler"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
	"log"
	"net/http"
)

type gophServer struct {
	storage        storage.Storage
	balanceHandler *handler.BalanceHandler
	orderHandler   *handler.OrderHandler
	userHandler    *handler.UserHandler
}

func NewGothServer(storage storage.Storage, serviceAddress string) *gophServer {
	ctx := context.Background()
	err := storage.Init(ctx)
	if err != nil {
		log.Printf("error while initializing storage: %v", err)
		return nil
	}
	balanceService := core.NewBalanceService(storage)
	orderService := core.NewOrderService(storage, serviceAddress, ctx)
	userService := core.NewUserService(storage)

	balanceHandler := handler.BalanceHandler{Service: balanceService, UserService: userService}
	orderHandler := handler.OrderHandler{Service: orderService, UserService: userService}
	userHandler := handler.UserHandler{Service: userService}

	return &gophServer{storage: storage, balanceHandler: &balanceHandler, orderHandler: &orderHandler, userHandler: &userHandler}
}

func (gs *gophServer) ServerHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(auth.CheckAuth)

		r.Post("/api/user/orders", func(rw http.ResponseWriter, r *http.Request) {
			gs.orderHandler.PostNewOrderHandler(rw, r)
		})

		r.Get("/api/user/orders", func(rw http.ResponseWriter, r *http.Request) {
			gs.orderHandler.GetOrders(rw, r)
		})

		r.Get("/api/user/balance", func(rw http.ResponseWriter, r *http.Request) {
			gs.balanceHandler.GetBalance(rw, r)
		})

		r.Post("/api/user/balance/withdraw", func(rw http.ResponseWriter, r *http.Request) {
			gs.balanceHandler.PostBalanceWithdraw(rw, r)
		})

		r.Get("/api/user/balance/withdrawals", func(rw http.ResponseWriter, r *http.Request) {
			gs.balanceHandler.GetBalanceWithdrawals(rw, r)
		})
	})

	r.Post("/api/user/register", func(rw http.ResponseWriter, r *http.Request) {
		gs.userHandler.RegisterHandler(rw, r)
	})

	r.Post("/api/user/login", func(rw http.ResponseWriter, r *http.Request) {
		gs.userHandler.LoginHandler(rw, r)
	})

	return r
}
