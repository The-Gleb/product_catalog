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
	registerURL = "/api/v1/register"
)

type RegisterUsecase interface {
	Register(ctx context.Context, credentials entity.Credentials) (entity.Session, error)
}

type registerHandler struct {
	usecase     RegisterUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewRegisterHandler(usecase RegisterUsecase) *registerHandler {
	return &registerHandler{usecase: usecase}
}

func (h *registerHandler) AddToRouter(r *chi.Mux) {
	r.Route(registerURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.Register)
	})

}

func (h *registerHandler) Middlewares(md ...func(http.Handler) http.Handler) *registerHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *registerHandler) Register(w http.ResponseWriter, r *http.Request) {

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

	s, err := h.usecase.Register(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {
		case errors.ErrAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
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
