package repository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
)

type authRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) portrepository.AuthRepository {
	return &authRepository{DB: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *domain.User) error {
	model := domainToModel(user)

	err := r.DB.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := User{}
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, gormx.MapError(err)
	}
	return user.ToDomain(), nil
}
