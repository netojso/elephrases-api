package service_test

import (
	"context"
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"github.com/stretchr/testify/mock"
)

type flashcardRepoMock struct{ mock.Mock }

func (m *flashcardRepoMock) FindByDeckID(ctx context.Context, deckID string) ([]*domain.Flashcard, error) {
	args := m.Called(ctx, deckID)
	return args.Get(0).([]*domain.Flashcard), args.Error(1)
}

func (m *flashcardRepoMock) FindStudyQueue(ctx context.Context, deckID string, dailyLimit int, now time.Time) (*domain.Study, error) {
	args := m.Called(ctx, deckID, dailyLimit, now)
	study, _ := args.Get(0).(*domain.Study)
	return study, args.Error(1)
}

func (m *flashcardRepoMock) FindAll(ctx context.Context, options *portrepository.Options) ([]*domain.Flashcard, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]*domain.Flashcard), args.Error(1)
}

func (m *flashcardRepoMock) FindByID(ctx context.Context, id string) (*domain.Flashcard, error) {
	args := m.Called(ctx, id)
	flashcard, _ := args.Get(0).(*domain.Flashcard)
	return flashcard, args.Error(1)
}

func (m *flashcardRepoMock) FindByIDForUpdate(ctx context.Context, id string) (*domain.Flashcard, error) {
	args := m.Called(ctx, id)
	flashcard, _ := args.Get(0).(*domain.Flashcard)
	return flashcard, args.Error(1)
}

func (m *flashcardRepoMock) Save(ctx context.Context, flashcard *domain.Flashcard) error {
	return m.Called(ctx, flashcard).Error(0)
}

func (m *flashcardRepoMock) Update(ctx context.Context, flashcard *domain.Flashcard) error {
	return m.Called(ctx, flashcard).Error(0)
}

func (m *flashcardRepoMock) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type reviewRepoMock struct{ mock.Mock }

func (m *reviewRepoMock) Save(ctx context.Context, review *domain.Review) error {
	return m.Called(ctx, review).Error(0)
}

type txCtxKey struct{}

// txManagerMock registra a chamada e executa fn com um ctx marcado, para os
// testes verificarem que os repositórios foram chamados dentro da transação.
type txManagerMock struct{ mock.Mock }

func (m *txManagerMock) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	m.Called(ctx)
	return fn(context.WithValue(ctx, txCtxKey{}, true))
}

func inTx(ctx context.Context) bool {
	v, _ := ctx.Value(txCtxKey{}).(bool)
	return v
}

type deckRepoMock struct{ mock.Mock }

func (m *deckRepoMock) FindAll(ctx context.Context) ([]*domain.Deck, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Deck), args.Error(1)
}

func (m *deckRepoMock) FindByID(ctx context.Context, id string) (*domain.Deck, error) {
	args := m.Called(ctx, id)
	deck, _ := args.Get(0).(*domain.Deck)
	return deck, args.Error(1)
}

func (m *deckRepoMock) Save(ctx context.Context, deck *domain.Deck) error {
	return m.Called(ctx, deck).Error(0)
}

func (m *deckRepoMock) Update(ctx context.Context, deck *domain.Deck) error {
	return m.Called(ctx, deck).Error(0)
}

func (m *deckRepoMock) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type userRepoMock struct{ mock.Mock }

func (m *userRepoMock) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *userRepoMock) Fetch(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *userRepoMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *userRepoMock) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *userRepoMock) UpdateUser(ctx context.Context, id string, user *domain.User) error {
	return m.Called(ctx, id, user).Error(0)
}

func (m *userRepoMock) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type authRepoMock struct{ mock.Mock }

func (m *authRepoMock) CreateUser(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *authRepoMock) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}
