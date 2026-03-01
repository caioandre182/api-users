package store

import (
	"context"

	"github.com/caioandre182/api-users/domain"
)

type UserStore interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, u domain.User) error
	Delete(ctx context.Context, id string) error
}
