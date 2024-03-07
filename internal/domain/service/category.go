package service

import (
	"context"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/usecase"
)

var _ usecase.CategoryService = new(categoryService)

type CategoryStorage interface {
	Add(ctx context.Context, Category entity.AddCategoryDTO) error
	GetAll(ctx context.Context) ([]entity.Category, error)
	UpdateName(ctx context.Context, category entity.UpdateCategoryNameDTO) error
	Delete(ctx context.Context, ID int64) error
}

type categoryService struct {
	storage CategoryStorage
}

func NewCategoryService(s CategoryStorage) *categoryService {
	return &categoryService{storage: s}
}

func (s *categoryService) Add(ctx context.Context, Category entity.AddCategoryDTO) error {
	return s.storage.Add(ctx, Category)
}

func (s *categoryService) GetAll(ctx context.Context) ([]entity.Category, error) {
	return s.storage.GetAll(ctx)
}

func (s *categoryService) UpdateName(ctx context.Context, category entity.UpdateCategoryNameDTO) error {
	return s.storage.UpdateName(ctx, category)
}

func (s *categoryService) Delete(ctx context.Context, ID int64) error {
	return s.storage.Delete(ctx, ID)
}
