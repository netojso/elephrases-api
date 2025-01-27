package usecase

import (
	"github.com/netojso/elephrases-api/domain"
	"github.com/netojso/elephrases-api/internal/token_util"
)

type refreshTokenUsecase struct {
	userRepository domain.UserRepository
}

func NewRefreshTokenUsecase(userRepository domain.UserRepository) domain.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		userRepository: userRepository,
	}
}

func (rtu *refreshTokenUsecase) GetUserByID(email string) (domain.User, error) {

	return rtu.userRepository.GetByID(email)
}

func (rtu *refreshTokenUsecase) CreateAccessToken(user domain.User, secret string, expiry int) (accessToken string, err error) {
	return token_util.CreateAccessToken(user, secret, expiry)
}

func (rtu *refreshTokenUsecase) CreateRefreshToken(user domain.User, secret string, expiry int) (refreshToken string, err error) {
	return token_util.CreateRefreshToken(user, secret, expiry)
}

func (rtu *refreshTokenUsecase) ExtractIDFromToken(requestToken string, secret string) (string, error) {
	return token_util.ExtractIDFromToken(requestToken, secret)
}
