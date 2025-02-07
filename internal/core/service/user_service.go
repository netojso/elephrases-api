package service

import (
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
)

type UserService struct {
	repo portrepository.UserRepository
}

func NewUserService(repo portrepository.UserRepository) portservice.UserService {
	return &UserService{repo: repo}
}

func (uu *UserService) Create(user *domain.User) error {
	return uu.repo.Create(user)
}

func (uu *UserService) GetByID(id string) (*domain.User, error) {
	return uu.repo.GetByID(id)
}

func (uu *UserService) GetByEmail(email string) (*domain.User, error) {
	return uu.repo.GetByEmail(email)
}

func (uu *UserService) Fetch() ([]*domain.User, error) {
	return uu.repo.Fetch()
}

func (uu *UserService) UpdateUser(id string, user *domain.User) error {

	_, err := uu.repo.GetByID(id)

	if err != nil {
		return err
	}

	return uu.repo.UpdateUser(id, user)
}

func (uu *UserService) DeleteUser(id string) error {

	_, err := uu.repo.GetByID(id)

	if err != nil {
		return err
	}

	return uu.repo.DeleteUser(id)
}
