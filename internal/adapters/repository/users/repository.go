package repository

import (
	"database/sql"
	"errors"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) portrepository.UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	model := domainToModel(user)

	err := r.DB.Create(&model).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Fetch() ([]*domain.User, error) {
	var users []*User
	err := r.DB.Find(&users).Error

	if err != nil {
		return nil, err
	}

	return modelToDomainList(users), nil

}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user User

	err := r.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	var user User
	err := r.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) UpdateUser(id string, user *domain.User) error {
	model := domainToModel(user)
	err := r.DB.Save(model).Error
	return err
}

func (r *UserRepository) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	err := r.DB.Exec(query, id).Error
	return err
}
