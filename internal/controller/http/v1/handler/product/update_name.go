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

const updateProductNameURL = "/api/v1/product/updateName"

type UpdateProductNameUsecase interface {
	UpdateName(ctx context.Context, product entity.UpdateProductNameDTO) error
}

type updateProductNameHandler struct {
	usecase     UpdateProductNameUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewUpdateProductNameHandler(usecase UpdateProductNameUsecase) *updateProductNameHandler {
	return &updateProductNameHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *updateProductNameHandler) AddToRouter(r *chi.Mux) {
	r.Route(updateProductNameURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

}

func (h *updateProductNameHandler) Middlewares(md ...func(http.Handler) http.Handler) *updateProductNameHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *updateProductNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var dto entity.UpdateProductNameDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("error decoding json request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.ProductID == 0 || dto.NewName == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateName(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())
		switch errors.Code(err) {
		case errors.ErrAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case errors.ErrNoDataFound:
			http.Error(w, err.Error(), http.StatusBadRequest)
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
