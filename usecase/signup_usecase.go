package usecase

import (
	"github.com/netojso/go-api-template/domain"
	"github.com/netojso/go-api-template/internal/token_util"
)

type signupUsecase struct {
	userRepository domain.UserRepository
}

func NewSignupUsecase(userRepository domain.UserRepository) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
	}
}

func (su *signupUsecase) Create(user domain.User) error {

	return su.userRepository.Create(user)
}

func (su *signupUsecase) GetUserByEmail(email string) (domain.User, error) {

	return su.userRepository.GetByEmail(email)
}

func (su *signupUsecase) CreateAccessToken(user domain.User, secret string, expiry int) (accessToken string, err error) {
	return token_util.CreateAccessToken(user, secret, expiry)
}

func (su *signupUsecase) CreateRefreshToken(user domain.User, secret string, expiry int) (refreshToken string, err error) {
	return token_util.CreateRefreshToken(user, secret, expiry)
}
