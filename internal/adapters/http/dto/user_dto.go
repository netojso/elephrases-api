package dto

import (
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type UpdateUserDTO struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type ResponseUserDTO struct {
	ID          pkg.UUID                `json:"id"`
	FullName    nullable.NullableString `json:"full_name"`
	Email       string                  `json:"email"`
	PhoneNumber nullable.NullableString `json:"phone_number"`
}
