package v1

import (
	"encoding/json"
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

func Test_addProductHandler_ServeHTTP(t *testing.T) {
	dto := entity.AddProductDTO{
		ProductName: "redmi",
		CategoryID:  1,
	}
	validRequestBody, err := json.Marshal(dto)
	require.NoError(t, err)

	invalidRequestBody, err := json.Marshal(entity.AddProductDTO{
		ProductName: "redmi",
	})
	require.NoError(t, err)

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockAddProductUsecase := mocks.NewMockAddProductUsecase(ctrl)
	loginHandler := NewAddProductHandler(mockAddProductUsecase)
	loginHandler.AddToRouter(r)
	server := httptest.NewServer(r)

	tests := []struct {
		name    string
		reqBody json.RawMessage
		code    int
		prepare func()
	}{
		{
			name:    "positive",
			reqBody: validRequestBody,
			code:    200,
			prepare: func() {
				mockAddProductUsecase.
					EXPECT().
					Add(gomock.Any(), gomock.Eq(dto)).
					Return(nil)
			},
		},
		{
			name:    "invalid body",
			reqBody: invalidRequestBody,
			code:    400,
			prepare: func() {},
		},
		{
			name:    "invalid body 2",
			reqBody: []byte("sdfasd"),
			code:    400,
			prepare: func() {},
		},
		{
			name:    "positive",
			reqBody: validRequestBody,
			code:    409,
			prepare: func() {
				mockAddProductUsecase.
					EXPECT().
					Add(gomock.Any(), gomock.Eq(dto)).
					Return(errors.NewDomainError(errors.ErrAlreadyExists, ""))
			},
		},
		{
			name:    "positive",
			reqBody: validRequestBody,
			code:    400,
			prepare: func() {
				mockAddProductUsecase.
					EXPECT().
					Add(gomock.Any(), gomock.Eq(dto)).
					Return(errors.NewDomainError(errors.ErrCategoryNotFound, ""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			resp, _ := v1.TestRequest(t, "", server, "POST", "/api/v1/product/add", tt.reqBody)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)
		})
	}
}
