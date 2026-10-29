package repository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := domainToModel(user)

	err := gormx.DB(ctx, r.DB).Create(&model).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Fetch(ctx context.Context) ([]*domain.User, error) {
	var users []*User
	err := gormx.DB(ctx, r.DB).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return modelToDomainList(users), nil

}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user User

	err := gormx.DB(ctx, r.DB).Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, gormx.MapError(err)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user User
	err := gormx.DB(ctx, r.DB).Where("id = ?", id).First(&user).Error

	if err != nil {
		return nil, gormx.MapError(err)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, user *domain.User) error {
	model := domainToModel(user)
	err := gormx.DB(ctx, r.DB).Save(model).Error
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = ?"
	return gormx.DeleteResult(gormx.DB(ctx, r.DB).Exec(query, id))
}
