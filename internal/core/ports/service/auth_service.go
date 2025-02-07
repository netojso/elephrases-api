package portservice

import "github.com/netojso/elephrases-api/internal/core/domain"

type AuthService interface {
	Login(email string, password string) (*domain.Session, error)
	Register(email string, password string) (*domain.Session, error)
}
