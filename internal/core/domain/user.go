package domain

import (
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type User struct {
	ID          pkg.UUID                `json:"id"`
	FullName    nullable.NullableString `json:"full_name"`
	Password    string                  `json:"-"`
	Email       string                  `json:"email"`
	PhoneNumber nullable.NullableString `json:"phone_number"`
}

func NewUser(email string, password string) *User {
	return &User{
		ID:       pkg.NewUUID(),
		Email:    email,
		Password: password,
	}
}

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           u.ID.String(),
		"full_name":    u.FullName.String,
		"email":        u.Email,
		"phone_number": u.PhoneNumber.String,
	}
}
