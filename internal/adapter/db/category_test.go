package db

import (
	"context"
	"testing"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func Test_categoryStorage_GetAll(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)

	_, err := client.Exec(
		context.Background(),
		`INSERT INTO category
			("id", "name")
		VALUES
			(1,'phone'),
			(2,'laptop'),
			(123, 'vacuum cleaner');`,
	)
	require.NoError(t, err)
	storage := NewCategoryStorage(client)

	tests := []struct {
		name      string
		want      []entity.Category
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name: "success",
			want: []entity.Category{
				{ID: 1, Name: "phone"},
				{ID: 2, Name: "laptop"},
				{ID: 123, Name: "vacuum cleaner"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories, err := storage.GetAll(context.Background())
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			require.ElementsMatch(t, tt.want, categories)

		})
	}
}

func Test_categoryStorage_Add(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)
	storage := NewCategoryStorage(client)

	tests := []struct {
		name      string
		dto       entity.AddCategoryDTO
		want      []entity.Category
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name: "success",
			dto: entity.AddCategoryDTO{
				Name: "phone",
			},
			wantErr: false,
		},
		{
			name: "already exists",
			dto: entity.AddCategoryDTO{
				Name: "phone",
			},
			wantErr:   true,
			errorCode: errors.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Add(context.Background(), tt.dto)
			if tt.wantErr {
				require.Equal(t, tt.errorCode, errors.Code(err))
				return
			}
			require.NoError(t, err)

			var name string
			row := client.QueryRow(
				context.Background(),
				`SELECT name FROM category
				WHERE name = $1;`,
				tt.dto.Name,
			)
			err = row.Scan(&name)
			require.NoError(t, err)

			require.Equal(t, tt.dto.Name, name)

		})
	}
}

func Test_categoryStorage_UpdateName(t *testing.T) {
	client := getTestClient(t)
	cleanTables(
		t, client,
		"product_category", "product", "category",
	)

	_, err := client.Exec(
		context.Background(),
		`INSERT INTO category
			("id", "name")
		VALUES
			(1,'phone'),
			(2,'laptop');`,
	)
	require.NoError(t, err)
	storage := NewCategoryStorage(client)

	tests := []struct {
		name      string
		dto       entity.UpdateCategoryNameDTO
		wantErr   bool
		errorCode errors.ErrorCode
	}{
		{
			name: "success",
			dto: entity.UpdateCategoryNameDTO{
				CategoryID: 1,
				NewName:    "earphones",
			},
			wantErr: false,
		},
		{
			name: "category id doesn't exist",
			dto: entity.UpdateCategoryNameDTO{
				CategoryID: 0,
				NewName:    "redmi",
			},
			wantErr:   true,
			errorCode: errors.ErrNoDataFound,
		},
		{
			name: "category name not unique",
			dto: entity.UpdateCategoryNameDTO{
				CategoryID: 1,
				NewName:    "laptop",
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
				`SELECT name FROM category
				WHERE id = $1;`,
				tt.dto.CategoryID,
			)
			err = row.Scan(&newName)
			require.NoError(t, err)

			require.Equal(t, tt.dto.NewName, newName)

		})
	}
}

func Test_categoryStorage_Delete(t *testing.T) {
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
	storage := NewCategoryStorage(client)

	tests := []struct {
		name       string
		idToDelete int64
		wantErr    bool
		errorCode  errors.ErrorCode
	}{
		{
			name:       "success, no references",
			idToDelete: 2,
			wantErr:    false,
		},
		{
			name:       "succes, with references",
			idToDelete: 123,
			wantErr:    false,
		},
		{
			name:       "id not found",
			idToDelete: 12,
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
				`SELECT * FROM category
				WHERE id = $1;`,
				tt.idToDelete,
			)
			err = row.Scan()
			require.Equal(t, pgx.ErrNoRows, err)

			row = client.QueryRow(
				context.Background(),
				`SELECT * FROM product_category
				WHERE category_id = $1;`,
				tt.idToDelete,
			)
			err = row.Scan()
			require.Equal(t, pgx.ErrNoRows, err)

		})
	}
}
