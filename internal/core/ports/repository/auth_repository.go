package portrepository

import "github.com/netojso/elephrases-api/internal/core/domain"

type AuthRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
}
