package user

import (
	"context"

	domain "reading-cats-api/internal/domain/user"
)

type Repository interface {
	FindByCognitoSub(ctx context.Context, sub domain.CognitoSub) (*domain.User, error)
	Insert(ctx context.Context, u *domain.User) error
}
