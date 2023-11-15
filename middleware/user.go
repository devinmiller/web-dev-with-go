package middleware

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/context"
	"github.com/devinmiller/web-dev-with-go/services"
)

type UserMiddleware struct {
	SessionService *services.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := umw.SessionService.CurrentUser(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
