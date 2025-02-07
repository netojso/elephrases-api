package repository

import (
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

func (r *authRepository) CreateUser(user *domain.User) error {
	model := domainToModel(user)

	err := r.DB.Create(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := User{}
	err := r.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}
