package handler

import (
	"encoding/json"
	"github.com/xbreathoflife/gophermart/internal/app/core"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	"io"
	"net/http"
)

type BalanceHandler struct {
	Service     *core.BalanceService
	UserService *core.UserService
}


func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := checkAuth(h.UserService, w, ctx)
	if sessionModel == nil {
		return
	}
	balance, err := h.Service.GetUsersBalance(ctx, sessionModel.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(balance)
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

func (h *BalanceHandler) PostBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
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
	bw := entities.BalanceWithdrawRequest{}
	if err := json.Unmarshal(b, &bw); err != nil {
		http.Error(w, "Error during parsing request json", http.StatusBadRequest)
		return
	}

	statusCode, err := h.Service.ProcessBalanceWithdraw(ctx, sessionModel.Login, bw)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(statusCode)
}

func (h *BalanceHandler) GetBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := checkAuth(h.UserService, w, ctx)
	if sessionModel == nil {
		return
	}

	orders, err := h.Service.GetWithdrawalsForUser(ctx, sessionModel.Login)
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
