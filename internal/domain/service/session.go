package service

import (
	"context"
	"time"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
	"github.com/The-Gleb/product_catalog/internal/domain/usecase"
	"github.com/The-Gleb/product_catalog/internal/errors"

	"github.com/google/uuid"
)

var _ usecase.SessionService = new(sessionService)

type SessionStorage interface {
	Create(ctx context.Context, session entity.Session) error
	GetByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error

	DeleteExpired(ctx context.Context) error
}

type sessionService struct {
	storage SessionStorage
}

func NewSessionService(s SessionStorage) *sessionService {
	return &sessionService{storage: s}
}

func (ss *sessionService) GetByToken(ctx context.Context, token string) (entity.Session, error) {

	err := ss.storage.DeleteExpired(ctx)
	if err != nil {
		return entity.Session{}, err
	}

	session, err := ss.storage.GetByToken(ctx, token)
	if err != nil {
		return entity.Session{}, err
	}
	if session.IsExpired() {
		return session, errors.NewDomainError(errors.ErrSessionExpired, "")
	}
	return session, nil
}

func (ss *sessionService) Create(ctx context.Context, userID int64) (entity.Session, error) {

	newSession := entity.Session{
		Token:  uuid.NewString(),
		UserID: userID,
		Expiry: time.Now().Add(24 * time.Hour),
	}

	for {
		err := ss.storage.Create(ctx, newSession)
		if errors.Code(err) == errors.ErrAlreadyExists {
			newSession.Token = uuid.NewString()
			continue
		}
		if err != nil {
			return entity.Session{}, err
		}
		break
	}
	return newSession, nil
}

func (ss *sessionService) Delete(ctx context.Context, token string) error {
	err := ss.storage.Delete(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
