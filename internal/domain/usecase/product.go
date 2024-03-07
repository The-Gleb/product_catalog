package usecase

import (
	"context"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
)

type productUsecase struct {
	productService ProductService
}

func NewProductUsecase(s ProductService) *productUsecase {
	return &productUsecase{
		productService: s,
	}
}

func (s *productUsecase) Add(ctx context.Context, product entity.AddProductDTO) error {
	return s.productService.Add(ctx, product)
}

func (s *productUsecase) GetByCategory(ctx context.Context, categoryID int64) ([]entity.ProductCategoryListItem, error) {
	return s.productService.GetByCategory(ctx, categoryID)
}

func (s *productUsecase) UpdateName(ctx context.Context, product entity.UpdateProductNameDTO) error {
	return s.productService.UpdateName(ctx, product)
}

func (s *productUsecase) UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error {
	return s.productService.UpdateCategory(ctx, product)
}

func (s *productUsecase) Delete(ctx context.Context, ID int64) error {
	return s.productService.Delete(ctx, ID)
}
