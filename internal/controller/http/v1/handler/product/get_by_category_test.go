package v1

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	v1 "github.com/The-Gleb/product_catalog/internal/controller/http/v1/handler"
	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_getProductByCategoryHandler_ServeHTTP_Success(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGetByCategoryUsecase := mocks.NewMockGetProductsByCategoryUsecase(ctrl)
	handler := NewGetProductsByCategoryHandler(mockGetByCategoryUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	id := int64(1)
	stringID := strconv.FormatInt(id, 10)

	mockGetByCategoryUsecase.EXPECT().GetByCategory(gomock.Any(), id).
		Return([]entity.ProductCategoryListItem{
			{ID: 1, Name: "iphone"},
		}, nil)

	resp, _ := v1.TestRequest(t, "", server, "GET", "/api/v1/product/get/"+stringID, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_getProductByCategoryHandler_ServeHTTP_InvalidID(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGetByCategoryUsecase := mocks.NewMockGetProductsByCategoryUsecase(ctrl)
	handler := NewGetProductsByCategoryHandler(mockGetByCategoryUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	resp, _ := v1.TestRequest(t, "", server, "GET", "/api/v1/product/get/adc", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_getProductByCategoryHandler_ServeHTTP_IDNotFound(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGetByCategoryUsecase := mocks.NewMockGetProductsByCategoryUsecase(ctrl)
	handler := NewGetProductsByCategoryHandler(mockGetByCategoryUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	id := int64(1)
	stringID := strconv.FormatInt(id, 10)

	mockGetByCategoryUsecase.EXPECT().GetByCategory(gomock.Any(), id).
		Return(nil, errors.NewDomainError(errors.ErrCategoryNotFound, ""))

	resp, _ := v1.TestRequest(t, "", server, "GET", "/api/v1/product/get/"+stringID, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_getProductByCategoryHandler_ServeHTTP_UnexpectedError(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGetByCategoryUsecase := mocks.NewMockGetProductsByCategoryUsecase(ctrl)
	handler := NewGetProductsByCategoryHandler(mockGetByCategoryUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	id := int64(1)
	stringID := strconv.FormatInt(id, 10)

	mockGetByCategoryUsecase.EXPECT().GetByCategory(gomock.Any(), id).
		Return(nil, errors.NewDomainError(errors.ErrDB, ""))

	resp, _ := v1.TestRequest(t, "", server, "GET", "/api/v1/product/get/"+stringID, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
