package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/The-Gleb/product_catalog/internal/controller/http/v1/handler"
	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_getAllCategoriesHandler_ServeHTTP_Success(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGetAllCategoriesUsecase := mocks.NewMockGetAllCategoriesUsecase(ctrl)
	handler := NewGetAllCategoriesHandler(mockGetAllCategoriesUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	expectedBody, err := json.Marshal([]entity.Category{
		{ID: 1, Name: "laptop"},
	})
	require.NoError(t, err)

	mockGetAllCategoriesUsecase.EXPECT().GetAll(gomock.Any()).
		Return([]entity.Category{
			{ID: 1, Name: "laptop"},
		}, nil)

	resp, body := v1.TestRequest(t, "", server, "GET", "/api/v1/category/getAll", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, string(expectedBody), body)
}

func Test_getAllCategoriesHandler_ServeHTTP_UnexpectedError(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGetAllCategoriesUsecase := mocks.NewMockGetAllCategoriesUsecase(ctrl)
	handler := NewGetAllCategoriesHandler(mockGetAllCategoriesUsecase)
	handler.AddToRouter(r)
	server := httptest.NewServer(r)

	mockGetAllCategoriesUsecase.EXPECT().GetAll(gomock.Any()).
		Return(nil, errors.NewDomainError(errors.ErrDB, ""))

	resp, _ := v1.TestRequest(t, "", server, "GET", "/api/v1/category/getAll", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
