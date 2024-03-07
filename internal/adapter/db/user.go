package db

import (
	"context"
	"log/slog"

	stdErrors "errors"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/service"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/pkg/client/postgresql"
	"github.com/jackc/pgx/v5"
)

var _ service.UserStorage = new(userStorage)

type userStorage struct {
	client postgresql.Client
}

func NewUserStorage(client postgresql.Client) *userStorage {
	return &userStorage{
		client: client,
	}
}

func (us *userStorage) Create(ctx context.Context, user entity.User) (entity.User, error) {

	row := us.client.QueryRow(
		ctx,
		`INSERT INTO "user"
			(login, password)
		VALUES
			($1,$2)
		ON CONFLICT DO NOTHING
		RETURNING *;`,
		user.Login, user.Password,
	)

	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		slog.Error("error adding user to db",
			"error", err,
		)
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "error adding user to db")
	}

	return user, nil
}

func (us *userStorage) GetByLogin(ctx context.Context, login string) (entity.User, error) {

	row := us.client.QueryRow(
		ctx,
		`SELECT * FROM "user"
		WHERE login = $1;`,
		login,
	)

	var user entity.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		slog.Error("error getting user from db",
			"error", err,
		)
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "")
		}
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "error getting user from db")
	}

	return user, nil

}

func (us *userStorage) GetByID(ctx context.Context, ID int64) (entity.User, error) {

	row := us.client.QueryRow(
		ctx,
		`SELECT * FROM "user"
		WHERE id = $1;`,
		ID,
	)

	var user entity.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		slog.Error("error getting user from db",
			"error", err,
		)
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "")
		}
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "error getting user from db")
	}

	return user, nil

}
