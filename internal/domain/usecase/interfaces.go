package usecase

import (
	"context"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
)

type SessionService interface {
	Create(ctx context.Context, userID int64) (entity.Session, error)
	GetByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, ID int64) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
}

type ProductService interface {
	Add(ctx context.Context, products entity.AddProductDTO) error
	GetByCategory(ctx context.Context, categoryID int64) ([]entity.ProductCategoryListItem, error)
	UpdateName(ctx context.Context, product entity.UpdateProductNameDTO) error
	UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error
	Delete(ctx context.Context, ID int64) error
}

type CategoryService interface {
	Add(ctx context.Context, Category entity.AddCategoryDTO) error
	GetAll(ctx context.Context) ([]entity.Category, error)
	UpdateName(ctx context.Context, category entity.UpdateCategoryNameDTO) error
	Delete(ctx context.Context, ID int64) error
}
