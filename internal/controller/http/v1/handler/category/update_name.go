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

const updateCategoryNameURL = "/api/v1/category/updateName"

type UpdateCategoryNameUsecase interface {
	UpdateName(ctx context.Context, category entity.UpdateCategoryNameDTO) error
}

type updateCategoryNameHandler struct {
	usecase     UpdateCategoryNameUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewUpdateCategoryNameHandler(usecase UpdateCategoryNameUsecase) *updateCategoryNameHandler {
	return &updateCategoryNameHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *updateCategoryNameHandler) AddToRouter(r *chi.Mux) {
	r.Route(updateCategoryNameURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

}

func (h *updateCategoryNameHandler) Middlewares(md ...func(http.Handler) http.Handler) *updateCategoryNameHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *updateCategoryNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var dto entity.UpdateCategoryNameDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("error decoding json request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.CategoryID == 0 || dto.NewName == "" {
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
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
