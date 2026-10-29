package portservice

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	Fetch(ctx context.Context) ([]*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}
