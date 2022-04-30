package auth

import (
	"context"
	"github.com/xbreathoflife/gophermart/internal/app/core"
	"net/http"
)

const CookieName = "authorization"
const CtxKey = ContextKey("session")

type ContextKey string


func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			session, err := core.Decrypt(cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), CtxKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}

