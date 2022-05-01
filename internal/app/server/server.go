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
	storage      storage.Storage
	handlers     *handler.Handler
}

func NewGothServer(storage storage.Storage, serviceAddress string) *gophServer {
	ctx := context.Background()
	err := storage.Init(ctx)
	if err != nil {
		log.Printf("error while initializing storage: %v", err)
		return nil
	}
	service := core.NewLoyaltyService(storage, serviceAddress, ctx)
	handlers := handler.Handler{Service: service}

	return &gophServer{storage: storage, handlers: &handlers}
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
			gs.handlers.PostNewOrderHandler(rw, r)
		})

		r.Get("/api/user/orders", func(rw http.ResponseWriter, r *http.Request) {
			gs.handlers.GetOrders(rw, r)
		})

		r.Get("/api/user/balance", func(rw http.ResponseWriter, r *http.Request) {
			gs.handlers.GetBalance(rw, r)
		})

		r.Post("/api/user/balance/withdraw", func(rw http.ResponseWriter, r *http.Request) {
			gs.handlers.PostBalanceWithdraw(rw, r)
		})

		r.Get("/api/user/balance/withdrawals", func(rw http.ResponseWriter, r *http.Request) {
			gs.handlers.GetBalanceWithdrawals(rw, r)
		})
	})

	r.Post("/api/user/register", func(rw http.ResponseWriter, r *http.Request) {
		gs.handlers.RegisterHandler(rw, r)
	})

	r.Post("/api/user/login", func(rw http.ResponseWriter, r *http.Request) {
		gs.handlers.LoginHandler(rw, r)
	})

	return r
}
