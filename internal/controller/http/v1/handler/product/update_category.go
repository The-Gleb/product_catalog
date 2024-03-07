package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

const updateProductCategoryURL = "/api/v1/product/updateCategory"

type UpdateProductCategoryUsecase interface {
	UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error
}

type updateProductCategoryHandler struct {
	usecase     UpdateProductCategoryUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewUpdateProductCategoryHandler(usecase UpdateProductCategoryUsecase) *updateProductCategoryHandler {
	return &updateProductCategoryHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *updateProductCategoryHandler) AddToRouter(r *chi.Mux) {
	r.Route(updateProductCategoryURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

}

func (h *updateProductCategoryHandler) Middlewares(md ...func(http.Handler) http.Handler) *updateProductCategoryHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *updateProductCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var dto entity.UpdateProductCategoryDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("error decoding json request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if dto.NewCategoryID == 0 || dto.OldCategoryID == 0 || dto.ProductID == 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateCategory(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
