package db

import (
	"context"
	"testing"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/stretchr/testify/require"
)

func Test_userStorage_Create(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"user",
	)
	userStorage := NewUserStorage(client)

	tests := []struct {
		name      string
		user      entity.User
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name:    "success",
			user:    entity.User{Login: "login1", Password: "passsword1"},
			wantErr: false,
		},
		{
			name:      "login already exists",
			user:      entity.User{Login: "login1", Password: "passsword1"},
			wantErr:   true,
			errorCode: errors.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userStorage.Create(context.Background(), tt.user)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}

			require.Equal(t, tt.user.Login, got.Login)
			require.Equal(t, tt.user.Password, got.Password)

		})
	}
}

func Test_userStorage_GetByLogin(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"user",
	)
	userStorage := NewUserStorage(client)
	user, err := userStorage.Create(
		context.Background(),
		entity.User{Login: "login1", Password: "passsword1"},
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		user      entity.User
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name:    "success",
			user:    user,
			wantErr: false,
		},
		{
			name:      "user doesn't exist",
			user:      entity.User{Login: "nouser", Password: "passsword1"},
			wantErr:   true,
			errorCode: errors.ErrNoDataFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userStorage.GetByLogin(context.Background(), tt.user.Login)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}

			require.Equal(t, tt.user, got)

		})
	}
}

func Test_userStorage_GetByID(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"user",
	)
	userStorage := NewUserStorage(client)
	user, err := userStorage.Create(
		context.Background(),
		entity.User{Login: "login1", Password: "passsword1"},
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		user      entity.User
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name:    "success",
			user:    user,
			wantErr: false,
		},
		{
			name:      "user doesn't exist",
			user:      entity.User{ID: 0, Login: "nouser", Password: "passsword1"},
			wantErr:   true,
			errorCode: errors.ErrNoDataFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userStorage.GetByID(context.Background(), tt.user.ID)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}

			require.Equal(t, tt.user, got)

		})
	}
}
