package service

import (
	"errors"

	"github.com/netojso/elephrases-api/config"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg"
)

type authService struct {
	repo portrepository.AuthRepository
	env  *config.Env
}

func NewAuthService(repo portrepository.AuthRepository, env *config.Env) portservice.AuthService {
	return &authService{repo: repo, env: env}
}

func (s *authService) Login(email string, password string) (*domain.Session, error) {
	user, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	valid := pkg.CompareHashAndPassword(user.Password, password)

	if !valid {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := pkg.CreateAccessToken(user.ToMap(), s.env.AccessTokenSecret, s.env.AccessTokenExpiryHour)

	if err != nil {
		return nil, err
	}

	refreshToken, err := pkg.CreateRefreshToken(user.ToMap(), s.env.RefreshTokenSecret, s.env.RefreshTokenExpiryHour)

	if err != nil {
		return nil, err
	}

	return &domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Register(email string, password string) (*domain.Session, error) {

	_, err := s.repo.GetUserByEmail(email)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	user := domain.NewUser(email, password)

	hash_password, err := pkg.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	user.Password = hash_password

	err = s.repo.CreateUser(user)

	if err != nil {
		return nil, err
	}

	accessToken, err := pkg.CreateAccessToken(user.ToMap(), s.env.AccessTokenSecret, s.env.AccessTokenExpiryHour)

	if err != nil {
		return nil, err
	}

	refreshToken, err := pkg.CreateRefreshToken(user.ToMap(), s.env.RefreshTokenSecret, s.env.RefreshTokenExpiryHour)

	if err != nil {
		return nil, err
	}

	return &domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
