package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/netojso/elephrases-api/config"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/internal/core/service"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testEnv = &config.Env{AccessTokenSecret: "access", RefreshTokenSecret: "refresh"}

func newAuthService(t *testing.T) (*authRepoMock, portservice.AuthService) {
	t.Helper()
	repo := &authRepoMock{}
	t.Cleanup(func() { repo.AssertExpectations(t) })
	return repo, service.NewAuthService(repo, testEnv)
}

func userWithPassword(t *testing.T, password string) *domain.User {
	t.Helper()
	hash, err := pkg.HashPassword(password)
	require.NoError(t, err)
	return domain.NewUser("a@x.com", hash)
}

func TestAuthService_Login(t *testing.T) {
	t.Run("credenciais válidas geram tokens", func(t *testing.T) {
		repo, svc := newAuthService(t)
		repo.On("GetUserByEmail", mock.Anything, "a@x.com").Return(userWithPassword(t, "secret"), nil)

		session, err := svc.Login(context.Background(), "a@x.com", "secret")

		require.NoError(t, err)
		assert.NotEmpty(t, session.AccessToken)
		assert.NotEmpty(t, session.RefreshToken)
	})

	t.Run("senha errada", func(t *testing.T) {
		repo, svc := newAuthService(t)
		repo.On("GetUserByEmail", mock.Anything, "a@x.com").Return(userWithPassword(t, "secret"), nil)

		_, err := svc.Login(context.Background(), "a@x.com", "wrong")

		assert.EqualError(t, err, "invalid credentials")
	})

	t.Run("e-mail inexistente responde igual a senha errada", func(t *testing.T) {
		repo, svc := newAuthService(t)
		repo.On("GetUserByEmail", mock.Anything, "zz@x.com").Return(nil, pkg.ErrNotFound)

		_, err := svc.Login(context.Background(), "zz@x.com", "secret")

		assert.EqualError(t, err, "invalid credentials")
	})
}

func TestAuthService_Register(t *testing.T) {
	t.Run("cria usuário com senha em hash", func(t *testing.T) {
		repo, svc := newAuthService(t)
		repo.On("GetUserByEmail", mock.Anything, "new@x.com").Return(nil, pkg.ErrNotFound)
		repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == "new@x.com" && u.Password != "secret" && pkg.CompareHashAndPassword(u.Password, "secret")
		})).Return(nil)

		session, err := svc.Register(context.Background(), "new@x.com", "secret")

		require.NoError(t, err)
		assert.NotEmpty(t, session.AccessToken)
	})

	t.Run("e-mail já cadastrado", func(t *testing.T) {
		repo, svc := newAuthService(t)
		repo.On("GetUserByEmail", mock.Anything, "a@x.com").Return(userWithPassword(t, "x"), nil)

		_, err := svc.Register(context.Background(), "a@x.com", "secret")

		assert.EqualError(t, err, "user already exists")
		repo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	})

	t.Run("erro do banco na consulta não cria usuário", func(t *testing.T) {
		repo, svc := newAuthService(t)
		boom := errors.New("connection refused")
		repo.On("GetUserByEmail", mock.Anything, "a@x.com").Return(nil, boom)

		_, err := svc.Register(context.Background(), "a@x.com", "secret")

		assert.ErrorIs(t, err, boom)
		repo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	})
}
