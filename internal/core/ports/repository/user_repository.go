package portrepository

import "github.com/netojso/elephrases-api/internal/core/domain"

type UserRepository interface {
	Create(user *domain.User) error
	Fetch() ([]*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	UpdateUser(id string, user *domain.User) error
	DeleteUser(id string) error
}
