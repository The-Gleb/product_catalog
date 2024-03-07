package v1

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/go-chi/chi/v5"
)

const deleteCategoryURL = "/api/v1/category/delete/{id}"

type DeleteCategoryUsecase interface {
	Delete(ctx context.Context, ID int64) error
}

type deleteCategoryHandler struct {
	usecase     DeleteCategoryUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewDeleteCategoryHandler(usecase DeleteCategoryUsecase) *deleteCategoryHandler {
	return &deleteCategoryHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *deleteCategoryHandler) AddToRouter(r *chi.Mux) {
	r.Route(deleteCategoryURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

}

func (h *deleteCategoryHandler) Middlewares(md ...func(http.Handler) http.Handler) *deleteCategoryHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *deleteCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	stringID := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		slog.Error("error parsing id from param to int64", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.usecase.Delete(r.Context(), ID)
	if err != nil {
		slog.Error(err.Error())
		switch errors.Code(err) {
		case errors.ErrNoDataFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
