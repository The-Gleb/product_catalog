package db

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func Test_sessionStorage_Delete(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"session", "user",
	)

	userStorage := NewUserStorage(client)
	user, err := userStorage.Create(
		context.Background(),
		entity.User{Login: "login1", Password: "password1"},
	)
	require.NoError(t, err)

	sessionStorage := NewSessionStorage(client)
	err = sessionStorage.Create(context.Background(), entity.Session{Token: "1", UserID: user.ID, Expiry: time.Now()})
	require.NoError(t, err)

	tests := []struct {
		name          string
		tokenToDelete string
		wantErr       bool
		errorCode     errors.ErrorCode
	}{
		{
			name:          "success",
			tokenToDelete: "1",
			wantErr:       false,
		},
		{
			name:          "token doesn't exist",
			tokenToDelete: "123",
			wantErr:       true,
			errorCode:     errors.ErrNoDataFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sessionStorage.Delete(context.Background(), tt.tokenToDelete)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			row := client.QueryRow(
				context.Background(),
				`SELECT * FROM session
				WHERE token = $1;`,
				tt.tokenToDelete,
			)
			err = row.Scan()
			require.Equal(t, pgx.ErrNoRows, err)

		})
	}
}

func Test_sessionStorage_DeleteExpired(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"session", "user",
	)
	userStorage := NewUserStorage(client)
	user, err := userStorage.Create(
		context.Background(),
		entity.User{Login: "login1", Password: "password1"},
	)
	require.NoError(t, err)

	sessionStorage := NewSessionStorage(client)
	sessionStorage.Create(context.Background(), entity.Session{Token: "1", UserID: user.ID, Expiry: time.Now()})
	sessionStorage.Create(context.Background(), entity.Session{Token: "2", UserID: user.ID, Expiry: time.Now()})
	sessionStorage.Create(context.Background(), entity.Session{Token: "3", UserID: user.ID, Expiry: time.Now()})
	time.Sleep(time.Second)

	err = sessionStorage.DeleteExpired(context.Background())
	require.NoError(t, err)

	rows, err := client.Query(
		context.Background(),
		`SELECT token FROM "session";`,
	)
	require.NoError(t, err)
	tokens, err := pgx.CollectRows[string](rows, func(row pgx.CollectableRow) (string, error) {
		var token string
		err := row.Scan(&token)
		return token, err
	})
	require.NoError(t, err)
	slog.Debug("tokens", "slice", tokens)
	require.Empty(t, tokens)

}
