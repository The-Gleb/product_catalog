package service

import (
	"context"
	"time"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/usecase"
)

var _ usecase.ProductService = new(productService)

type ProductStorage interface {
	Add(ctx context.Context, products entity.AddProductDTO) error
	AddOrUpdateProduct(ctx context.Context, products ...entity.AddOrUpdateProductDTO) error
	GetByCategory(ctx context.Context, categoryID int64) ([]entity.ProductCategoryListItem, error)
	UpdateName(ctx context.Context, product entity.UpdateProductNameDTO) error
	UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error
	Delete(ctx context.Context, ID int64) error
}

type ProductClient interface {
	GetNewProducts(ctx context.Context, offset int) ([]entity.AddOrUpdateProductDTO, error)
}

type productService struct {
	storage        ProductStorage
	client         ProductClient
	updateInterval time.Duration
}

func NewProductService(s ProductStorage, c ProductClient, interval time.Duration) *productService {
	return &productService{
		storage:        s,
		client:         c,
		updateInterval: interval,
	}
}

func (s *productService) Add(ctx context.Context, product entity.AddProductDTO) error {
	return s.storage.Add(ctx, product)
}

func (s *productService) GetByCategory(ctx context.Context, categoryID int64) ([]entity.ProductCategoryListItem, error) {
	return s.storage.GetByCategory(ctx, categoryID)
}

func (s *productService) UpdateName(ctx context.Context, product entity.UpdateProductNameDTO) error {
	return s.storage.UpdateName(ctx, product)
}

func (s *productService) UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error {
	return s.storage.UpdateCategory(ctx, product)
}

func (s *productService) Delete(ctx context.Context, ID int64) error {
	return s.storage.Delete(ctx, ID)
}

func (s *productService) CheckNewProducts(ctx context.Context) error { // TODO: handle errors

	ticker := time.NewTicker(s.updateInterval)
	offset := 10

	newProducts, err := s.client.GetNewProducts(ctx, 0)
	if err != nil {
		return err
	}
	err = s.storage.AddOrUpdateProduct(ctx, newProducts...)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ticker.C:
			newProducts, err := s.client.GetNewProducts(ctx, offset)
			if err != nil {
				return err
			}
			err = s.storage.AddOrUpdateProduct(ctx, newProducts...)
			if err != nil {
				return err
			}
			offset = offset + 10
		case <-ctx.Done():
			return nil
		}
	}

}
