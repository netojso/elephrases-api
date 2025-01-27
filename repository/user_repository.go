package repository

import (
	"database/sql"
	"errors"

	"github.com/netojso/go-api-template/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user domain.User) error {
	query := "INSERT INTO users (id, email, password) VALUES (?, ?, ?)"
	err := r.DB.Exec(query, user.ID, user.Email, user.Password).Error
	return err
}

func (r *UserRepository) Fetch() ([]domain.User, error) {
	query := `SELECT
  id,
  full_name,
  password,
  email,
  phone_number
	FROM users`

	rows, err := r.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domain.User, 0)

	for rows.Next() {
		user := new(domain.User)
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Password,
			&user.Email,
			&user.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil

}

func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
	query := `SELECT
  id,
  full_name,
  password,
  email,
  phone_number
	FROM users
 	WHERE email = ?`

	row := r.DB.Raw(query, email).Row()

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}
	return *user, nil
}

func (r *UserRepository) GetByID(id string) (domain.User, error) {
	// query := "SELECT id, full_name, email, password FROM users WHERE id = ?"
	query := `SELECT
  id,
  full_name,
  password,
  email,
  phone_number
	FROM users 
	WHERE id = ?`

	row := r.DB.Raw(query, id).Row()
	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return *user, nil
}

func (r *UserRepository) UpdateUser(id string, user domain.User) error {

	query := `
		UPDATE users 
		SET 
			full_name = COALESCE(?, full_name),
			phone_number = COALESCE(?, phone_number),
		WHERE id = ?
	`
	err := r.DB.Exec(query, user.FullName, user.PhoneNumber, id).Error
	return err
}

func (r *UserRepository) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	err := r.DB.Exec(query, id).Error
	return err
}
