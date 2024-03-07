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

const addProductURL = "/api/v1/product/add"

type AddProductUsecase interface {
	Add(ctx context.Context, product entity.AddProductDTO) error
}

type addProductHandler struct {
	usecase     AddProductUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewAddProductHandler(usecase AddProductUsecase) *addProductHandler {
	return &addProductHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *addProductHandler) AddToRouter(r *chi.Mux) {
	r.Route(addProductURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

}

func (h *addProductHandler) Middlewares(md ...func(http.Handler) http.Handler) *addProductHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *addProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dto entity.AddProductDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("error decoding json request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.CategoryID == 0 || dto.ProductName == "" {
		http.Error(w, "empty name or category id", http.StatusBadRequest)
		return
	}

	err = h.usecase.Add(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())
		switch errors.Code(err) {
		case errors.ErrAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case errors.ErrCategoryNotFound:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
