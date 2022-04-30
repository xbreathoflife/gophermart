package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/xbreathoflife/gophermart/internal/app/auth"
	"github.com/xbreathoflife/gophermart/internal/app/core"
	"github.com/xbreathoflife/gophermart/internal/app/entities"
	er "github.com/xbreathoflife/gophermart/internal/app/errors"
	"io"
	"net/http"
)

type Handler struct {
	Service *core.LoyaltyService
}

func (h *Handler) processCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, *string) {
	_, err := r.Cookie(auth.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil
	}
	uuid := core.GenerateUUID()
	encryptedUUID, err := core.Encrypt(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil
	}
	return &http.Cookie{Name: auth.CookieName, Value: encryptedUUID}, &uuid
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := entities.LoginRequest{}
	if err := json.Unmarshal(b, &user); err != nil {
		http.Error(w, "Error during parsing request json", http.StatusBadRequest)
		return
	}
	if user.Login == "" || user.Password == "" {
		http.Error(w, "Password or login empty", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.Service.CheckUserExists(ctx, user)
	if err != nil {
		var de *er.DuplicateError
		if errors.As(err, &de) {
			http.Error(w, "Username is taken", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//cookie
	newCookie, uuid := h.processCookie(w, r)
	if newCookie == nil {
		return
	}
	err = h.Service.InsertNewUser(ctx, entities.UserModel{Login: user.Login, PasswordHash: user.Password, Session: *uuid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, newCookie)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := entities.LoginRequest{}
	if err := json.Unmarshal(b, &user); err != nil {
		http.Error(w, "Error during parsing request json", http.StatusBadRequest)
		return
	}
	if user.Login == "" || user.Password == "" {
		http.Error(w, "Password or login empty", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.Service.CheckUserCredentials(ctx, user)
	if err != nil {
		var ce *er.WrongDataError
		if errors.As(err, &ce) {
			http.Error(w, "Wrong username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//cookie
	newCookie, uuid := h.processCookie(w, r)
	if newCookie == nil {
		return
	}
	err = h.Service.UpdateUserSession(ctx, entities.UserSessionModel{Login: user.Login, Session: *uuid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, newCookie)
}

func (h *Handler) checkAuth(w http.ResponseWriter, ctx context.Context) *entities.UserSessionModel {
	session := ctx.Value(auth.CtxKey).(string)
	sessionModel, err := h.Service.GetUserBySession(ctx, session)
	if err != nil {
		var ce *er.WrongDataError
		if errors.As(err, &ce) {
			http.Error(w, "Wrong username or password", http.StatusUnauthorized)
			return nil
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return sessionModel
}

func (h *Handler) PostNewOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := h.checkAuth(w, ctx)
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

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := h.checkAuth(w, ctx)
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

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := h.checkAuth(w, ctx)
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

func (h *Handler) PostBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := h.checkAuth(w, ctx)
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

func (h *Handler) GetBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionModel := h.checkAuth(w, ctx)
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
