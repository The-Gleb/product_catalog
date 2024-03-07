package db

import (
	"context"
	stdErrors "errors"
	"fmt"
	"log/slog"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/service"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/pkg/client/postgresql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ service.CategoryStorage = new(categoryStorage)

type categoryStorage struct {
	client postgresql.Client
}

func NewCategoryStorage(client postgresql.Client) *categoryStorage {
	return &categoryStorage{
		client: client,
	}
}

func (s *categoryStorage) Add(ctx context.Context, category entity.AddCategoryDTO) error {
	query := fmt.Sprintf(
		`INSERT INTO "category"
			(name)
		VALUES
			('%s');`,
		category.Name,
	)
	slog.Debug(query)
	c, err := s.client.Exec(
		ctx,
		query,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		slog.Error("error inserting into category",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if c.RowsAffected() == 0 {
		return errors.NewDomainError(errors.ErrAlreadyExists, "")
	}
	return nil
}

func (s *categoryStorage) GetAll(ctx context.Context) ([]entity.Category, error) {

	rows, err := s.client.Query(
		ctx,
		`SELECT * FROM category;`,
	)
	if err != nil {
		slog.Error("error selcting from category",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}
	defer rows.Close()

	cats, err := pgx.CollectRows[entity.Category](
		rows, func(row pgx.CollectableRow) (entity.Category, error) {
			var cat entity.Category
			err := row.Scan(&cat.ID, &cat.Name)
			return cat, err
		},
	)
	if err != nil {
		slog.Error("error collecting rows",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}

	return cats, nil

}

func (s *categoryStorage) UpdateName(ctx context.Context, category entity.UpdateCategoryNameDTO) error {

	c, err := s.client.Exec(
		ctx,
		`UPDATE category
		SET name = $1
		WHERE id = $2;`,
		category.NewName,
		category.CategoryID,
	)
	if err != nil {
		slog.Error("error updating category name",
			"error", err,
		)
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if c.RowsAffected() == 0 {
		slog.Error("no rows affected, not unique category name",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrNoDataFound, "")
	}

	return nil

}

func (s *categoryStorage) Delete(ctx context.Context, ID int64) error {

	c, err := s.client.Exec(
		ctx,
		`DELETE FROM category
		WHERE id = $1;`,
		ID,
	)
	if err != nil {
		slog.Error("error deleting from category",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	if c.RowsAffected() == 0 {
		slog.Error("error deleting from category",
			"error", "no rows affected, id not found",
		)
		return errors.NewDomainError(errors.ErrNoDataFound, "")
	}
	return nil

}
