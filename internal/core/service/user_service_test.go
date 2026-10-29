package service_test

import (
	"context"
	"testing"

	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/internal/core/service"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_UpdateUser_NotFoundSkipsUpdate(t *testing.T) {
	repo := &userRepoMock{}
	t.Cleanup(func() { repo.AssertExpectations(t) })
	repo.On("GetByID", mock.Anything, "missing").Return(nil, pkg.ErrNotFound)

	err := service.NewUserService(repo).UpdateUser(context.Background(), "missing", &domain.User{})

	assert.ErrorIs(t, err, pkg.ErrNotFound)
	repo.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_DeleteUser(t *testing.T) {
	repo := &userRepoMock{}
	t.Cleanup(func() { repo.AssertExpectations(t) })
	user := domain.NewUser("a@x.com", "hash")
	repo.On("GetByID", mock.Anything, user.ID.String()).Return(user, nil)
	repo.On("DeleteUser", mock.Anything, user.ID.String()).Return(nil)

	require.NoError(t, service.NewUserService(repo).DeleteUser(context.Background(), user.ID.String()))
}
