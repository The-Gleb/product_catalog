package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

const getAllCategoriesURL = "/api/v1/category/getAll"

type GetAllCategoriesUsecase interface {
	GetAll(ctx context.Context) ([]entity.Category, error)
}

type getAllCategoriesHandler struct {
	usecase     GetAllCategoriesUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewGetAllCategoriesHandler(usecase GetAllCategoriesUsecase) *getAllCategoriesHandler {
	return &getAllCategoriesHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *getAllCategoriesHandler) AddToRouter(r *chi.Mux) {
	r.Route(getAllCategoriesURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

}

func (h *getAllCategoriesHandler) Middlewares(md ...func(http.Handler) http.Handler) *getAllCategoriesHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *getAllCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	categories, err := h.usecase.GetAll(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(categories)
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
