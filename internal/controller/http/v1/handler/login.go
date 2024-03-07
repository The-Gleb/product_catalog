package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/go-chi/chi/v5"
)

const (
	loginURL = "/api/v1/login"
)

type LoginUsecase interface {
	Login(ctx context.Context, credentials entity.Credentials) (entity.Session, error)
}

type loginHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     LoginUsecase
}

func NewLoginHandler(usecase LoginUsecase) *loginHandler {
	return &loginHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *loginHandler) AddToRouter(r *chi.Mux) {

	r.Route(loginURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.Login)
	})
}

func (h *loginHandler) Middlewares(md ...func(http.Handler) http.Handler) *loginHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *loginHandler) Login(w http.ResponseWriter, r *http.Request) {

	var dto entity.Credentials
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("error parsing json request body to dto", "error", err)
		http.Error(w, "error parsing json request body to dto", http.StatusBadRequest)
		return
	}

	if dto.Login == "" || dto.Password == "" {
		http.Error(w, "login and password should not be empty", http.StatusBadRequest)
		return
	}

	s, err := h.usecase.Login(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {
		case errors.ErrUnauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	c := http.Cookie{
		Name:    "sessionToken",
		Value:   s.Token,
		Expires: s.Expiry,
		Path:    "/",
	}

	http.SetCookie(w, &c)

	w.WriteHeader(http.StatusOK)

}
