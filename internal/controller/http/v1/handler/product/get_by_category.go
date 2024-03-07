package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/go-chi/chi/v5"
)

const getProductsByCategoryURL = "/api/v1/product/get/{categoryId}"

type GetProductsByCategoryUsecase interface {
	GetByCategory(ctx context.Context, categoryID int64) ([]entity.ProductCategoryListItem, error)
}

type getProductsByCategoryHandler struct {
	usecase     GetProductsByCategoryUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewGetProductsByCategoryHandler(usecase GetProductsByCategoryUsecase) *getProductsByCategoryHandler {
	return &getProductsByCategoryHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *getProductsByCategoryHandler) AddToRouter(r *chi.Mux) {
	r.Route(getProductsByCategoryURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

}

func (h *getProductsByCategoryHandler) Middlewares(md ...func(http.Handler) http.Handler) *getProductsByCategoryHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *getProductsByCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	stringID := chi.URLParam(r, "categoryId")
	ID, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		slog.Error("error parsing id from param to int64", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	products, err := h.usecase.GetByCategory(r.Context(), ID)
	if err != nil {
		slog.Error(err.Error())
		switch errors.Code(err) {
		case errors.ErrCategoryNotFound:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	body, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
