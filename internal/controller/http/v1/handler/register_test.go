package v1

import (
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_registerHandler_Register(t *testing.T) {
	validRegisterReqBody, err := json.Marshal(entity.Credentials{
		Login:    "login1",
		Password: "password1",
	})
	require.NoError(t, err)

	invalidRegisterReqBody, err := json.Marshal(entity.Credentials{
		Login: "login1",
	})
	require.NoError(t, err)

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)

	mockRegisterUsecase := mocks.NewMockRegisterUsecase(ctrl)
	registerHandler := NewRegisterHandler(mockRegisterUsecase)
	registerHandler.AddToRouter(r)
	server := httptest.NewServer(r)

	tests := []struct {
		name    string
		reqBody json.RawMessage
		code    int
		prepare func()
	}{
		{
			name:    "positive",
			reqBody: validRegisterReqBody,
			code:    200,
			prepare: func() {
				mockRegisterUsecase.
					EXPECT().
					Register(gomock.Any(), gomock.Eq(entity.Credentials{
						Login:    "login1",
						Password: "password1",
					})).
					Return(entity.Session{
						UserID: 1,
						Token:  "123",
						Expiry: time.Now().Add(time.Hour),
					}, nil)
			},
		},
		{
			name:    "negative, body with no username",
			reqBody: invalidRegisterReqBody,
			code:    400,
			prepare: func() {
			},
		},
		{
			name:    "negative, invalid body",
			reqBody: []byte("domasldkfjdf"),
			code:    400,
			prepare: func() {
			},
		},
		{
			name:    "negative, wrong login/password",
			reqBody: validRegisterReqBody,
			code:    409,
			prepare: func() {
				mockRegisterUsecase.
					EXPECT().
					Register(gomock.Any(), gomock.Eq(entity.Credentials{
						Login:    "login1",
						Password: "password1",
					})).
					Return(entity.Session{}, errors.NewDomainError(errors.ErrAlreadyExists, ""))
			},
		},
		{
			name:    "negative, some db err",
			reqBody: validRegisterReqBody,
			code:    500,
			prepare: func() {
				mockRegisterUsecase.
					EXPECT().
					Register(gomock.Any(), gomock.Eq(entity.Credentials{
						Login:    "login1",
						Password: "password1",
					})).
					Return(entity.Session{}, errors.NewDomainError(errors.ErrDB, ""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.prepare()

			resp, _ := TestRequest(t, "", server, "POST", "/api/v1/register", tt.reqBody)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)

			if tt.code != 200 {
				return
			}

			cookies := resp.Cookies()
			require.NotEqual(t, "0", len(cookies))
			require.Equal(t, "sessionToken", cookies[0].Name)
			require.NotEmpty(t, cookies[0].Value)
			slog.Debug("received cookie", "cookie", cookies[0])

		})
	}
}
