package service

import (
	"context"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/usecase"
)

var _ usecase.UserService = new(userService)

type UserStorage interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, ID int64) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
}

type userService struct {
	storage UserStorage
}

func NewUserService(s UserStorage) *userService {
	return &userService{storage: s}
}

func (us *userService) GetByID(ctx context.Context, ID int64) (entity.User, error) {
	return us.storage.GetByID(ctx, ID)
}

func (us *userService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	return us.storage.Create(ctx, user)
}

func (us *userService) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	return us.storage.GetByLogin(ctx, login)
}
