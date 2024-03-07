package v1

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	v1 "github.com/The-Gleb/product_catalog/internal/controller/http/v1/handler"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_deleteProductHandler_ServeHTTP_SuccessfulDeletion(t *testing.T) {
	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockDeleteProductUsecase := mocks.NewMockDeleteProductUsecase(ctrl)
	handler := NewDeleteProductHandler(mockDeleteProductUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	id := int64(1)
	stringID := strconv.FormatInt(id, 10)

	mockDeleteProductUsecase.EXPECT().Delete(gomock.Any(), id).Return(nil)

	resp, _ := v1.TestRequest(t, "", server, "POST", "/api/v1/product/delete/"+stringID, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_deleteProductHandler_ServeHTTP_InvalidID(t *testing.T) {
	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockDeleteProductUsecase := mocks.NewMockDeleteProductUsecase(ctrl)
	handler := NewDeleteProductHandler(mockDeleteProductUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	stringID := "invalid"

	resp, _ := v1.TestRequest(t, "", server, "POST", "/api/v1/product/delete/"+stringID, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_deleteProductHandler_ServeHTTP_UnknownError(t *testing.T) {
	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockDeleteProductUsecase := mocks.NewMockDeleteProductUsecase(ctrl)
	handler := NewDeleteProductHandler(mockDeleteProductUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	mockDeleteProductUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.NewDomainError(errors.ErrDB, ""))

	resp, _ := v1.TestRequest(t, "", server, "POST", "/api/v1/product/delete/1", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteProductHandler_ServeHTTP_ProductNotFound(t *testing.T) {
	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockDeleteProductUsecase := mocks.NewMockDeleteProductUsecase(ctrl)
	deleteHandler := NewDeleteProductHandler(mockDeleteProductUsecase)
	deleteHandler.AddToRouter(r)
	server := httptest.NewServer(r)

	mockDeleteProductUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.NewDomainError(errors.ErrNoDataFound, ""))

	resp, _ := v1.TestRequest(t, "", server, "POST", "/api/v1/product/delete/1", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
