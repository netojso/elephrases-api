package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	FullName    sql.NullString `gorm:"type:varchar(255)" json:"full_name"`
	Password    string         `gorm:"type:varchar(255)" json:"password"`
	Email       string         `gorm:"type:varchar(255);unique_index" json:"email"`
	PhoneNumber sql.NullString `gorm:"type:varchar(20)" json:"phone_number"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ToDomain() *domain.User {
	id, err := pkg.ParseUUID(u.ID.String())

	if err != nil {
		return nil
	}

	return &domain.User{
		ID:          id,
		FullName:    nullable.NewNullableString(u.FullName.String),
		Password:    u.Password,
		Email:       u.Email,
		PhoneNumber: nullable.NewNullableString(u.PhoneNumber.String),
	}
}

func domainToModel(user *domain.User) *User {
	return &User{
		ID:          user.ID.Value(),
		FullName:    user.FullName.NullString,
		Password:    user.Password,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber.NullString,
	}
}

func domainToModelList(users []*domain.User) []*User {
	var usersModel []*User

	for _, user := range users {
		usersModel = append(usersModel, domainToModel(user))
	}

	return usersModel
}

func modelToDomainList(users []*User) []*domain.User {
	var usersDomain []*domain.User

	for _, user := range users {
		usersDomain = append(usersDomain, user.ToDomain())
	}

	return usersDomain
}
