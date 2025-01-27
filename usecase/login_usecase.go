package usecase

import (
	"github.com/netojso/go-api-template/domain"
	"github.com/netojso/go-api-template/internal/token_util"
)

type loginUsecase struct {
	userRepository domain.UserRepository
}

func NewLoginUsecase(userRepository domain.UserRepository) domain.LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
	}
}

func (lu *loginUsecase) GetUserByEmail(email string) (domain.User, error) {

	return lu.userRepository.GetByEmail(email)
}

func (lu *loginUsecase) CreateAccessToken(user domain.User, secret string, expiry int) (accessToken string, err error) {
	return token_util.CreateAccessToken(user, secret, expiry)
}

func (lu *loginUsecase) CreateRefreshToken(user domain.User, secret string, expiry int) (refreshToken string, err error) {
	return token_util.CreateRefreshToken(user, secret, expiry)
}
