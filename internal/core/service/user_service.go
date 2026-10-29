package service

import (
	"context"

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

func (uu *UserService) Create(ctx context.Context, user *domain.User) error {
	return uu.repo.Create(ctx, user)
}

func (uu *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return uu.repo.GetByID(ctx, id)
}

func (uu *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uu.repo.GetByEmail(ctx, email)
}

func (uu *UserService) Fetch(ctx context.Context) ([]*domain.User, error) {
	return uu.repo.Fetch(ctx)
}

func (uu *UserService) UpdateUser(ctx context.Context, id string, user *domain.User) error {

	_, err := uu.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	return uu.repo.UpdateUser(ctx, id, user)
}

func (uu *UserService) DeleteUser(ctx context.Context, id string) error {

	_, err := uu.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	return uu.repo.DeleteUser(ctx, id)
}
