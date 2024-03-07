package db

import (
	"bytes"
	"context"
	stdErrors "errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/service"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/pkg/client/postgresql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ service.ProductStorage = new(productStorage)

type productStorage struct {
	client postgresql.Client
}

func NewProductStorage(client postgresql.Client) *productStorage {
	return &productStorage{
		client: client,
	}
}

func (ps *productStorage) AddOrUpdateProduct(ctx context.Context, products ...entity.AddOrUpdateProductDTO) error {
	if len(products) == 0 {
		slog.Error("products slice is emty")
		return nil
	}

	tx, err := ps.client.Begin(ctx)
	if err != nil {

		slog.Error("error beginnig transaction",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	defer tx.Rollback(ctx)

	productBuffer := bytes.Buffer{}
	categoryBuffer := bytes.Buffer{}

	productMap := make(map[string]entity.Product, len(products))
	categoryMap := make(map[string]int64, 0)
	for _, product := range products {
		_, err := productBuffer.WriteString(fmt.Sprintf("('%s'),", product.ProductName))
		if err != nil {
			slog.Error("error writing string to buffer",
				"error", err,
			)
			return errors.NewDomainError(errors.ErrDB, "")
		}
		productMap[product.ProductName] = entity.Product{
			Name: product.ProductName,
			Category: entity.Category{
				Name: product.CategoryName,
			},
		}
		if _, ok := categoryMap[product.CategoryName]; !ok {
			categoryMap[product.CategoryName] = 0
			_, err := categoryBuffer.WriteString(fmt.Sprintf("('%s'),", product.CategoryName))
			if err != nil {
				slog.Error("error writing string to buffer",
					"error", err,
				)
				return errors.NewDomainError(errors.ErrDB, "")
			}
		}
	}

	productNames := strings.TrimSuffix(productBuffer.String(), ",")
	categoryNames := strings.TrimSuffix(categoryBuffer.String(), ",")

	query := fmt.Sprintf(`
		INSERT INTO
			product ("name")
		VALUES %s
		ON CONFLICT(name)
		DO UPDATE SET
			name=EXCLUDED.name
		RETURNING
			id, name;
	`, productNames)

	rows, err := tx.Query(
		ctx,
		query,
	)
	if err != nil {
		slog.Error("error inserting products",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			slog.Error("error scanning from row",
				"error", err,
			)
			return errors.NewDomainError(errors.ErrDB, "")
		}
		p := productMap[name]
		p.ID = id
		productMap[name] = p
	}
	if err := rows.Err(); err != nil {
		slog.Error("error scanning from row",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	query = fmt.Sprintf(`
		INSERT INTO
			category ("name")
		VALUES %s
		ON CONFLICT(name)
		DO UPDATE SET
			name=EXCLUDED.name
		RETURNING
			id, name;
	`, categoryNames)

	rows, err = tx.Query(
		ctx,
		query,
	)
	if err != nil {
		slog.Error("error inserting category",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			slog.Error("error scanning from row",
				"error", err,
			)
			return errors.NewDomainError(errors.ErrDB, "")
		}
		categoryMap[name] = id
	}

	productCategoryBuffer := bytes.Buffer{}

	for _, product := range productMap {
		_, err := productCategoryBuffer.WriteString(
			fmt.Sprintf("(%d,%d),", product.ID, categoryMap[product.Category.Name]),
		)
		if err != nil {
			slog.Error("error writing to buffer",
				"error", err,
			)
			return errors.NewDomainError(errors.ErrDB, "")
		}
	}

	productCategoryString := strings.TrimSuffix(productCategoryBuffer.String(), ",")

	query = fmt.Sprintf(`
		INSERT INTO
			product_category ("product_id", "category_id")
		VALUES %s
		ON CONFLICT DO NOTHING;
	`, productCategoryString)

	slog.Debug(query)

	_, err = tx.Exec(
		ctx,
		query,
	)
	if err != nil {
		slog.Error("error inserting into product_category",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	defer rows.Close()

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error commiting transaction",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	return nil

}

func (ps *productStorage) Add(ctx context.Context, product entity.AddProductDTO) error {
	tx, err := ps.client.Begin(ctx)
	if err != nil {
		slog.Error("error beginnig transaction",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	defer tx.Rollback(ctx)

	exists, err := ps.categoryExists(ctx, product.CategoryID, tx)
	if err != nil {
		slog.Error("error chekcing if category exists",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if !exists {
		return errors.NewDomainError(errors.ErrCategoryNotFound, "")
	}

	row := tx.QueryRow(
		ctx,
		`INSERT INTO product
			("name")
		VALUES
			($1)
		ON CONFLICT DO NOTHING
		RETURNING id;`,
		product.ProductName,
	)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		slog.Error("error scanning from row",
			"error", err,
		)
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		return err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO product_category
			(product_id, category_id)
		VALUES
			($1,$2);`,
		id, product.CategoryID,
	)
	if err != nil {
		slog.Error("error inserting into product_category",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error commiting transaction",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	return nil

}

func (ps *productStorage) GetByCategory(ctx context.Context, categoryID int64) ([]entity.ProductCategoryListItem, error) {
	tx, err := ps.client.Begin(ctx)
	if err != nil {
		slog.Error("error beginnig transaction",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}
	defer tx.Rollback(ctx)

	exists, err := ps.categoryExists(ctx, categoryID, tx)
	if err != nil {
		slog.Error("error chekcing if category exists",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}
	if !exists {
		return nil, errors.NewDomainError(errors.ErrCategoryNotFound, "")
	}

	productIDs, err := ps.getProductIDsByCategory(ctx, categoryID, tx)
	if err != nil {
		slog.Error("error product id by category",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}

	query := fmt.Sprintf(
		`SELECT * FROM product
		WHERE id IN (%s);`,
		strings.Join(productIDs, ","),
	)

	rows, err := tx.Query(
		ctx,
		query,
	)
	if err != nil {
		slog.Error("error selecting from product table",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}

	list, err := pgx.CollectRows[entity.ProductCategoryListItem](
		rows, func(row pgx.CollectableRow) (entity.ProductCategoryListItem, error) {
			var product entity.ProductCategoryListItem
			err := row.Scan(&product.ID, &product.Name)
			return product, err
		},
	)
	if err != nil {
		slog.Error("error scanning rows",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error commiting transaction",
			"error", err,
		)
		return nil, errors.NewDomainError(errors.ErrDB, "")
	}

	return list, nil

}

func (ps *productStorage) UpdateName(ctx context.Context, product entity.UpdateProductNameDTO) error {
	c, err := ps.client.Exec(
		ctx,
		`UPDATE product
		SET name = $1
		WHERE id = $2;`,
		product.NewName,
		product.ProductID,
	)

	if err != nil {
		slog.Error("error product id by category",
			"error", err,
		)
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if c.RowsAffected() == 0 {
		slog.Error("no rows affected, product id not found")

		return errors.NewDomainError(errors.ErrNoDataFound, "")
	}
	return nil
}

func (ps *productStorage) UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error {
	tx, err := ps.client.Begin(ctx)
	if err != nil {
		slog.Error("error beginnig transaction",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	defer tx.Rollback(ctx)

	c, err := tx.Exec(
		ctx,
		`UPDATE product_category
		SET category_id = $2
		WHERE product_id = $1 AND category_id = $3;`,
		product.ProductID,
		product.NewCategoryID,
		product.OldCategoryID,
	)
	if err != nil {
		slog.Error("error product category",
			"error", err,
		)
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrAlreadyExists, "")
		}
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return errors.NewDomainError(errors.ErrCategoryNotFound, "")
		}
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if c.RowsAffected() == 0 {
		slog.Error("product with that category not found")
		return errors.NewDomainError(errors.ErrNoDataFound, "no rows affected!!!")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error commiting transaction",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}

	return nil

}

func (ps *productStorage) Delete(ctx context.Context, ID int64) error {
	c, err := ps.client.Exec(
		ctx,
		`DELETE FROM product CASCADE
		WHERE id = $1
		;`,
		ID,
	)
	if err != nil {
		slog.Error("error deleting from products",
			"error", err,
		)
		return errors.NewDomainError(errors.ErrDB, "")
	}
	if c.RowsAffected() == 0 {
		slog.Error("error deleting from products, id not found")
		return errors.NewDomainError(errors.ErrNoDataFound, "")
	}

	return nil
}

func (ps *productStorage) getProductIDsByCategory(ctx context.Context, categoryID int64, tx pgx.Tx) ([]string, error) {
	rows, err := tx.Query(
		ctx,
		`SELECT
			product_id
		FROM
			product_category
		WHERE
			category_id = $1;`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productIDs, err := pgx.CollectRows[string](rows, func(row pgx.CollectableRow) (string, error) {
		var id string
		err := row.Scan(&id)
		return id, err
	})
	if err != nil {
		return nil, err
	}

	return productIDs, nil

}

func (ps *productStorage) categoryExists(ctx context.Context, categoryID int64, tx pgx.Tx) (bool, error) {
	row := tx.QueryRow(
		ctx,
		`SELECT CASE WHEN EXISTS (
			SELECT * FROM category
			WHERE id = $1
		)
		THEN TRUE
		ELSE FALSE END;`,
		categoryID,
	)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
