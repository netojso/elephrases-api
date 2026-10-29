package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/internal/core/service"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeckService_GetAll_ComputesStats(t *testing.T) {
	repo := &deckRepoMock{}
	t.Cleanup(func() { repo.AssertExpectations(t) })

	deck := domain.NewDeck("d", "", "c", "v")
	recent := nullable.NewNullableTime(time.Now().Add(-time.Hour))
	deck.Flashcards = []domain.Flashcard{
		{State: domain.StateNew},
		{State: domain.StateNew},
		{State: domain.StateLearning, LastReviewAt: recent},
		{State: domain.StateReview, LastReviewAt: recent},
		{State: domain.StateReview},
	}
	repo.On("FindAll", mock.Anything).Return([]*domain.Deck{deck}, nil)

	decks, err := service.NewDeckUsecase(repo).GetAll(context.Background())

	require.NoError(t, err)
	require.Len(t, decks, 1)
	assert.Equal(t, domain.Stats{NewCards: 2, LearningCards: 1, ReviewCards: 2, TotalCards: 5}, decks[0].Stats)
}

func TestDeckService_Delete_PropagatesNotFound(t *testing.T) {
	repo := &deckRepoMock{}
	t.Cleanup(func() { repo.AssertExpectations(t) })
	repo.On("Delete", mock.Anything, "missing").Return(pkg.ErrNotFound)

	err := service.NewDeckUsecase(repo).Delete(context.Background(), "missing")

	assert.ErrorIs(t, err, pkg.ErrNotFound)
}
