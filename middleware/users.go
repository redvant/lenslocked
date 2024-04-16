package middleware

import (
	"net/http"

	"github.com/redvant/lenslocked/context"
	"github.com/redvant/lenslocked/cookies"
	"github.com/redvant/lenslocked/models"
)

type Users struct {
	SessionService *models.SessionService
}

func (u Users) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := cookies.ReadCookie(r, cookies.CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := u.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
