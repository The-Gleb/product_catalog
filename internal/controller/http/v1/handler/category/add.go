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

const addCategoryURL = "/api/v1/category/add"

type AddCategoryUsecase interface {
	Add(ctx context.Context, category entity.AddCategoryDTO) error
}

type addCategoryHandler struct {
	usecase     AddCategoryUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewAddCategoryHandler(usecase AddCategoryUsecase) *addCategoryHandler {
	return &addCategoryHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *addCategoryHandler) AddToRouter(r *chi.Mux) {
	r.Route(addCategoryURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

}

func (h *addCategoryHandler) Middlewares(md ...func(http.Handler) http.Handler) *addCategoryHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *addCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dto entity.AddCategoryDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("error decoding json request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.Name == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.usecase.Add(r.Context(), dto)
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

	w.WriteHeader(http.StatusOK)
}
