package usecase

import (
	"context"
	"log/slog"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type registerUsecase struct {
	userService    UserService
	sessionService SessionService
}

func NewRegisterUsecase(us UserService, ss SessionService) *registerUsecase {
	return &registerUsecase{us, ss}
}

func (uc *registerUsecase) Register(ctx context.Context, credentials entity.Credentials) (entity.Session, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.Session{}, err
	}
	credentials.Password = string(hashedPassword)

	slog.Debug("credentials", "struct", credentials)

	user, err := uc.userService.Create(ctx, entity.User{
		Login:    credentials.Login,
		Password: credentials.Password,
	})
	if err != nil {
		return entity.Session{}, err
	}

	newSession, err := uc.sessionService.Create(ctx, user.ID)
	if err != nil {
		return entity.Session{}, err
	}

	return newSession, nil
}
