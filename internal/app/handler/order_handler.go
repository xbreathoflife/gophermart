package handler

import (
	"encoding/json"
	"github.com/xbreathoflife/gophermart/internal/app/core"
	"io"
	"net/http"
)

type OrderHandler struct {
	Service     *core.OrderService
	UserService *core.UserService
}


func (h *OrderHandler) PostNewOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := checkAuth(h.UserService, w, ctx)
	if sessionModel == nil {
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	orderNum := string(b)
	statusCode, err := h.Service.CreateNewOrder(ctx, sessionModel.Login, orderNum)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(statusCode)
}


func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := checkAuth(h.UserService, w, ctx)
	if sessionModel == nil {
		return
	}

	orders, err := h.Service.GetOrdersForUser(ctx, sessionModel.Login)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	js, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Error during building response json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
