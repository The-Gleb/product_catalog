package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/product_catalog/internal/errors"
)

type Key string

type AuthUsecase interface {
	Auth(ctx context.Context, token string) (int64, error)
}

type authMiddleWare struct {
	usecase AuthUsecase
}

func NewAuthMiddleware(usecase AuthUsecase) *authMiddleWare {
	return &authMiddleWare{usecase}
}

func (m *authMiddleWare) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("auth middleware working")

		c, err := r.Cookie("sessionToken")
		if err != nil {
			slog.Error("error getting cookie", "error", err)
			http.Error(w, string(errors.ErrUnauthorized), http.StatusUnauthorized)
			return
		}
		slog.Debug("Cookie is", "cookie", c.Value)

		userID, err := m.usecase.Auth(r.Context(), c.Value)
		if err != nil {
			http.Error(w, string(errors.ErrUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), Key("userID"), userID)
		ctx = context.WithValue(ctx, Key("token"), c.Value)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
