package usecase

import (
	"context"
)

type authUsecase struct {
	sessionService SessionService
}

func NewAuthUsecase(ss SessionService) *authUsecase {
	return &authUsecase{ss}
}

func (uc *authUsecase) Auth(ctx context.Context, token string) (int64, error) {

	session, err := uc.sessionService.GetByToken(ctx, token)
	if err != nil {
		return 0, err
	}

	return session.UserID, nil
}
