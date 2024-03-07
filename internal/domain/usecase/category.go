package usecase

import (
	"context"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
)

type categoryUsecase struct {
	categoryService CategoryService
}

func NewCategoryUsecase(s CategoryService) *categoryUsecase {
	return &categoryUsecase{
		categoryService: s,
	}
}

func (s *categoryUsecase) Add(ctx context.Context, category entity.AddCategoryDTO) error {
	return s.categoryService.Add(ctx, category)
}

func (s *categoryUsecase) GetAll(ctx context.Context) ([]entity.Category, error) {
	return s.categoryService.GetAll(ctx)
}

func (s *categoryUsecase) UpdateName(ctx context.Context, category entity.UpdateCategoryNameDTO) error {
	return s.categoryService.UpdateName(ctx, category)
}

func (s *categoryUsecase) Delete(ctx context.Context, ID int64) error {
	return s.categoryService.Delete(ctx, ID)
}
