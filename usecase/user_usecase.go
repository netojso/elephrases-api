package usecase

import (
	"github.com/netojso/elephrases-api/domain"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository}
}

func (uu *userUsecase) Create(user domain.User) error {

	return uu.userRepository.Create(user)
}

func (uu *userUsecase) GetByID(id string) (domain.User, error) {

	return uu.userRepository.GetByID(id)
}

func (uu *userUsecase) GetByEmail(email string) (domain.User, error) {

	return uu.userRepository.GetByEmail(email)
}

func (uu *userUsecase) Fetch() ([]domain.User, error) {

	return uu.userRepository.Fetch()
}

func (uu *userUsecase) UpdateUser(id string, user domain.User) error {

	return uu.userRepository.UpdateUser(id, user)
}

func (uu *userUsecase) DeleteUser(id string) error {

	return uu.userRepository.DeleteUser(id)
}
