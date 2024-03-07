package db

import (
	"context"
	stdErrors "errors"
	"log/slog"
	"time"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/service"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/pkg/client/postgresql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ service.SessionStorage = new(sessionStorage)

type sessionStorage struct {
	client postgresql.Client
}

func NewSessionStorage(client postgresql.Client) *sessionStorage {
	return &sessionStorage{
		client: client,
	}
}

func (ss *sessionStorage) GetByToken(ctx context.Context, token string) (entity.Session, error) {

	row := ss.client.QueryRow(
		ctx,
		`SELECT * FROM session
		WHERE token = $1;`,
		token,
	)

	var session entity.Session
	err := row.Scan(&session.ID, &session.Token, &session.UserID, &session.Expiry)
	if err != nil {
		slog.Error("error getting session from db",
			"error", err,
		)
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.Session{}, errors.NewDomainError(errors.ErrNoDataFound, "")
		}
		return entity.Session{}, errors.NewDomainError(errors.ErrDB, "")
	}

	return session, nil

}

func (ss *sessionStorage) Create(ctx context.Context, session entity.Session) error {

	_, err := ss.client.Exec(
		ctx,
		`INSERT INTO session
			("token", "user_id", "expiry")
		VALUES
			($1,$2,$3);`,
		session.Token,
		session.UserID,
		session.Expiry,
	)

	if err != nil {
		slog.Error(err.Error())
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		slog.Error("error adding new session to db",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	return nil
}

func (ss *sessionStorage) Delete(ctx context.Context, token string) error {
	c, err := ss.client.Exec(
		ctx,
		`DELETE FROM session
		WHERE "token" = $1;`,
		token,
	)
	if err != nil {
		slog.Error("error deleting session with token",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if c.RowsAffected() == 0 {
		return errors.NewDomainError(errors.ErrNoDataFound, "no rows affected")
	}

	return nil
}

func (ss *sessionStorage) DeleteExpired(ctx context.Context) error {
	_, err := ss.client.Exec(
		ctx,
		`DELETE FROM session
		WHERE "expiry" < $1;`,
		time.Now(),
	)
	if err != nil {
		slog.Error("error deleting expired sessions",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	return nil
}
