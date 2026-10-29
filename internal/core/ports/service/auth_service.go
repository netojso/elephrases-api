package portservice

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*domain.Session, error)
	Register(ctx context.Context, email string, password string) (*domain.Session, error)
}
