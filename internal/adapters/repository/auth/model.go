package repository

import (
	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID               `gorm:"type:uuid;primary_key"`
	FullName    nullable.NullableString `gorm:"type:varchar(100)"`
	Password    string                  `gorm:"type:varchar(100);not null"`
	Email       string                  `gorm:"type:varchar(100);uniqueIndex;not null"`
	PhoneNumber nullable.NullableString `gorm:"type:varchar(15)"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func (user *User) ToDomain() *domain.User {

	uuid, err := pkg.ParseUUID(user.ID.String())

	if err != nil {
		return nil
	}

	return &domain.User{
		ID:          uuid,
		FullName:    user.FullName,
		Password:    user.Password,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
}

func domainToModel(user *domain.User) *User {
	return &User{
		ID:          user.ID.Value(),
		FullName:    user.FullName,
		Password:    user.Password,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
}
