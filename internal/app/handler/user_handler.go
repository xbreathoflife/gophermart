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

type UserHandler struct {
	Service *core.UserService
}


func checkAuth(service *core.UserService, w http.ResponseWriter, ctx context.Context) *entities.UserSessionModel {
	session := ctx.Value(auth.CtxKey).(string)
	sessionModel, err := service.GetUserBySession(ctx, session)
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


func (h *UserHandler) processCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, *string) {
	_, _ = r.Cookie(auth.CookieName)

	uuid := core.GenerateUUID()
	encryptedUUID, err := core.Encrypt(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil
	}
	return &http.Cookie{Name: auth.CookieName, Value: encryptedUUID}, &uuid
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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


func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
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