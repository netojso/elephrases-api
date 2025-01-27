package domain

import (
	"github.com/netojso/go-api-template/internal"
)

type User struct {
	ID          string                  `json:"id" validate:"required" gorm:"primaryKey"`
	FullName    internal.NullableString `json:"full_name" gorm:"type:varchar(255)"`
	Password    string                  `json:"password"`
	Email       string                  `json:"email" gorm:"uniqueIndex"`
	PhoneNumber internal.NullableString `json:"phone_number,omitempty" gorm:"type:varchar(20)"`
}

type UserRepository interface {
	Create(user User) error
	Fetch() ([]User, error)
	GetByEmail(email string) (User, error)
	GetByID(id string) (User, error)
	UpdateUser(id string, user User) error
	DeleteUser(id string) error
}

type UserUsecase interface {
	Create(user User) error
	Fetch() ([]User, error)
	GetByEmail(email string) (User, error)
	GetByID(id string) (User, error)
	UpdateUser(id string, user User) error
	DeleteUser(id string) error
}
