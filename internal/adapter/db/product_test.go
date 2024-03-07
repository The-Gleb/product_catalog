package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/The-Gleb/product_catalog/internal/config"
	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/The-Gleb/product_catalog/pkg/client/postgresql"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestClient(t *testing.T) postgresql.Client {
	ctx := context.Background()
	client, err := postgresql.NewClient(ctx, config.Database{
		Host:     "localhost",
		Port:     5434,
		Username: "catalog_db",
		Password: "catalog_db",
		DbName:   "catalog_db",
	})
	require.NoError(t, err)
	return client

}

func cleanTables(t *testing.T, client postgresql.Client, tableNames ...string) {
	for _, name := range tableNames {
		query := fmt.Sprintf("TRUNCATE TABLE \"%s\" CASCADE", name)
		_, err := client.Exec(
			context.Background(),
			query,
		)
		require.NoError(t, err)
	}

}

func Test_productStorage_AddOrUpdateProduct(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)
	storage := NewProductStorage(client)

	tests := []struct {
		name      string
		products  []entity.AddOrUpdateProductDTO
		result    []entity.AddOrUpdateProductDTO
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name: "first insert",
			products: []entity.AddOrUpdateProductDTO{
				{ProductName: "redmi", CategoryName: "phone"},
				{ProductName: "iphone", CategoryName: "phone"},
				{ProductName: "maibenben", CategoryName: "laptop"},
			},
			result: []entity.AddOrUpdateProductDTO{
				{ProductName: "redmi", CategoryName: "phone"},
				{ProductName: "iphone", CategoryName: "phone"},
				{ProductName: "maibenben", CategoryName: "laptop"},
			},
			wantErr: false,
		},
		{
			name: "success test",
			products: []entity.AddOrUpdateProductDTO{
				{ProductName: "redmi", CategoryName: "tablet"},
				{ProductName: "lenovo", CategoryName: "laptop"},
				{ProductName: "dyson", CategoryName: "vacuum cleaner"},
			},
			result: []entity.AddOrUpdateProductDTO{
				{ProductName: "redmi", CategoryName: "phone"},
				{ProductName: "redmi", CategoryName: "tablet"},
				{ProductName: "lenovo", CategoryName: "laptop"},
				{ProductName: "iphone", CategoryName: "phone"},
				{ProductName: "dyson", CategoryName: "vacuum cleaner"},
				{ProductName: "maibenben", CategoryName: "laptop"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.AddOrUpdateProduct(context.Background(), tt.products...)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			rows, err := client.Query(
				context.Background(),
				`SELECT product.name, category.name
				FROM product
				JOIN product_category ON product.id = product_category.product_id
				JOIN category ON product_category.category_id = category.id;`,
			)
			require.NoError(t, err)
			products, err := pgx.CollectRows[entity.AddOrUpdateProductDTO](
				rows, func(row pgx.CollectableRow) (entity.AddOrUpdateProductDTO, error) {
					var p entity.AddOrUpdateProductDTO
					err := row.Scan(&p.ProductName, &p.CategoryName)
					return p, err
				},
			)
			require.NoError(t, err)

			assert.ElementsMatch(t, tt.result, products)

		})
	}
}

func Test_productStorage_Add(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)
	row := client.QueryRow(
		context.Background(),
		`INSERT INTO category
			("name")
		VALUES
			('laptop')
		RETURNING id;`,
	)
	var catID int64
	err := row.Scan(&catID)
	require.NoError(t, err)
	storage := NewProductStorage(client)

	tests := []struct {
		name         string
		productToAdd entity.AddProductDTO
		result       entity.AddProductDTO
		wantErr      bool
		errorCode    errors.ErrorCode
	}{
		{
			name: "first insert",
			productToAdd: entity.AddProductDTO{
				ProductName: "redmi",
				CategoryID:  catID,
			},
			result: entity.AddProductDTO{
				ProductName: "redmi",
				CategoryID:  catID,
			},
			wantErr: false,
		},
		{
			name: "second insert duplicate",
			productToAdd: entity.AddProductDTO{
				ProductName: "redmi",
				CategoryID:  catID,
			},
			wantErr:   true,
			errorCode: errors.ErrAlreadyExists,
		},
		{
			name: "category doesn't exist",
			productToAdd: entity.AddProductDTO{
				ProductName: "iphone",
				CategoryID:  0,
			},
			wantErr:   true,
			errorCode: errors.ErrCategoryNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.Add(context.Background(), tt.productToAdd)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			row := client.QueryRow(
				context.Background(),
				`SELECT product.name, product_category.category_id
				FROM product
				JOIN product_category ON product.id = product_category.product_id
				WHERE product.name = $1 AND product_category.category_id = $2;`,
				tt.productToAdd.ProductName, tt.productToAdd.CategoryID,
			)
			require.NoError(t, err)

			var product entity.AddProductDTO
			err = row.Scan(&product.ProductName, &product.CategoryID)
			require.NoError(t, err)

			assert.Equal(t, tt.result, product)

		})
	}
}

func Test_productStorage_GetByCategory(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)
	storage := NewProductStorage(client)
	err := storage.AddOrUpdateProduct(
		context.Background(),
		[]entity.AddOrUpdateProductDTO{
			{ProductName: "redmi", CategoryName: "phone"},
			{ProductName: "lenovo", CategoryName: "laptop"},
			{ProductName: "iphone", CategoryName: "phone"},
			{ProductName: "dyson", CategoryName: "vacuum cleaner"},
			{ProductName: "maibenben", CategoryName: "laptop"},
		}...,
	)
	require.NoError(t, err)

	var phoneCategoryID int64
	row := client.QueryRow(
		context.Background(),
		`SELECT id FROM category
		WHERE name = 'phone';`,
	)
	err = row.Scan(&phoneCategoryID)
	require.NoError(t, err)

	var iphoneID int64
	row = client.QueryRow(
		context.Background(),
		`SELECT id FROM product
		WHERE name = 'iphone';`,
	)
	err = row.Scan(&iphoneID)
	require.NoError(t, err)

	var redmiID int64
	row = client.QueryRow(
		context.Background(),
		`SELECT id FROM product
		WHERE name = 'redmi';`,
	)
	err = row.Scan(&redmiID)
	require.NoError(t, err)

	tests := []struct {
		name       string
		categoryID int64
		result     []entity.ProductCategoryListItem
		wantErr    bool
		errorCode  errors.ErrorCode
	}{
		{
			name:       "existing category",
			categoryID: phoneCategoryID,
			result: []entity.ProductCategoryListItem{
				{ID: redmiID, Name: "redmi"},
				{ID: iphoneID, Name: "iphone"},
			},
			wantErr: false,
		},
		{
			name:       "non existing category",
			categoryID: 0,
			wantErr:    true,
			errorCode:  errors.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			products, err := storage.GetByCategory(context.Background(), tt.categoryID)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			assert.ElementsMatch(t, products, tt.result)

		})
	}
}

func Test_productStorage_UpdateName(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)

	_, err := client.Exec(
		context.Background(),
		`INSERT INTO product
			("id", "name")
		VALUES
			(1,'iphone'),
			(2,'redmi');`,
	)
	require.NoError(t, err)
	storage := NewProductStorage(client)

	tests := []struct {
		name      string
		dto       entity.UpdateProductNameDTO
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name: "success",
			dto: entity.UpdateProductNameDTO{
				ProductID: 1,
				NewName:   "galaxy",
			},
			wantErr: false,
		},
		{
			name: "product id doesn't exist",
			dto: entity.UpdateProductNameDTO{
				ProductID: 0,
				NewName:   "redmi",
			},
			wantErr:   true,
			errorCode: errors.ErrNoDataFound,
		},
		{
			name: "product name not unique",
			dto: entity.UpdateProductNameDTO{
				ProductID: 1,
				NewName:   "redmi",
			},
			wantErr:   true,
			errorCode: errors.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.UpdateName(context.Background(), tt.dto)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			var newName string
			row := client.QueryRow(
				context.Background(),
				`SELECT name FROM product
				WHERE id = $1;`,
				tt.dto.ProductID,
			)
			err = row.Scan(&newName)
			require.NoError(t, err)

			assert.Equal(t, tt.dto.NewName, newName)

		})
	}

}

func Test_productStorage_UpdateCategory(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)

	_, err := client.Exec(
		context.Background(),
		`INSERT INTO product
			("id", "name")
		VALUES
			(1,'iphone'),
			(2,'redmi');

		INSERT INTO category
			("id", "name")
		VALUES
			(1,'phone'),
			(2,'laptop'),
			(123, 'vacuum cleaner');
		INSERT INTO product_category
			("product_id", "category_id")
		VALUES
			(1,1),
			(2,1),
			(2,123);`,
	)
	require.NoError(t, err)
	storage := NewProductStorage(client)

	tests := []struct {
		name      string
		dto       entity.UpdateProductCategoryDTO
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name:    "success",
			dto:     entity.UpdateProductCategoryDTO{ProductID: 2, OldCategoryID: 1, NewCategoryID: 2},
			wantErr: false,
		},
		{
			name:      "product already belongs to category",
			dto:       entity.UpdateProductCategoryDTO{ProductID: 2, OldCategoryID: 2, NewCategoryID: 123},
			wantErr:   true,
			errorCode: errors.ErrAlreadyExists,
		},
		{
			name:      "product with that category not found",
			dto:       entity.UpdateProductCategoryDTO{ProductID: 2, OldCategoryID: 1, NewCategoryID: 123},
			wantErr:   true,
			errorCode: errors.ErrNoDataFound,
		},
		{
			name:      "category doesn't exist",
			dto:       entity.UpdateProductCategoryDTO{ProductID: 2, OldCategoryID: 2, NewCategoryID: 5},
			wantErr:   true,
			errorCode: errors.ErrCategoryNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.UpdateCategory(context.Background(), tt.dto)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			var newCategoryID int64
			row := client.QueryRow(
				context.Background(),
				`SELECT category_id FROM product_category
				WHERE product_id = $1 AND category_id = $2;`,
				tt.dto.ProductID, tt.dto.NewCategoryID,
			)
			err = row.Scan(&newCategoryID)
			require.NoError(t, err)

			assert.Equal(t, tt.dto.NewCategoryID, newCategoryID)
		})
	}
}

func Test_productStorage_Delete(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)

	_, err := client.Exec(
		context.Background(),
		`INSERT INTO product
			("id", "name")
		VALUES
			(1,'iphone'),
			(2,'redmi');
		INSERT INTO category
			("id", "name")
		VALUES
			(1,'phone'),
			(2,'laptop'),
			(123, 'vacuum cleaner');
		INSERT INTO product_category
			("product_id", "category_id")
		VALUES
			(1,1),
			(2,1),
			(2,123);`,
	)
	require.NoError(t, err)
	storage := NewProductStorage(client)

	tests := []struct {
		name       string
		idToDelete int64
		wantErr    bool
		errorCode  errors.ErrorCode
	}{
		{
			name:       "success",
			idToDelete: 2,
			wantErr:    false,
		},
		{
			name:       "product doesn't exist",
			idToDelete: 4,
			wantErr:    true,
			errorCode:  errors.ErrNoDataFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Delete(context.Background(), tt.idToDelete)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			row := client.QueryRow(
				context.Background(),
				`SELECT * FROM product
				WHERE id = $1;`,
				tt.idToDelete,
			)
			err = row.Scan()
			require.Equal(t, pgx.ErrNoRows, err)

		})
	}
}
