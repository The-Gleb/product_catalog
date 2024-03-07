package usecase

import (
	"context"
	"log/slog"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	userService    UserService
	sessionService SessionService
}

func NewLoginUsecase(us UserService, ss SessionService) *loginUsecase {
	return &loginUsecase{us, ss}
}

func (uc *loginUsecase) Login(ctx context.Context, credentials entity.Credentials) (entity.Session, error) {

	user, err := uc.userService.GetByLogin(ctx, credentials.Login)
	if err != nil {
		if errors.Code(err) == errors.ErrNoDataFound {
			return entity.Session{}, errors.NewDomainError(errors.ErrUnauthorized, "")
		}
		return entity.Session{}, err
	}

	slog.Debug("user", "struct", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		slog.Error("error comparing passwords", "error", err)
		return entity.Session{}, errors.NewDomainError(errors.ErrUnauthorized, "")
	}

	newSession, err := uc.sessionService.Create(ctx, user.ID)
	if err != nil {
		return entity.Session{}, err
	}

	return newSession, nil
}
