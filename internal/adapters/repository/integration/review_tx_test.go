//go:build integration

package integration

import (
	"context"
	"errors"
	"testing"

	flashcardrepository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
	reviewrepository "github.com/netojso/elephrases-api/internal/adapters/repository/reviews"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/internal/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errReviewInsert = errors.New("simulated reviews insert failure")

// failingReviewRepository simula falha no insert em reviews.
type failingReviewRepository struct{}

func (failingReviewRepository) Save(context.Context, *domain.Review) error {
	return errReviewInsert
}

func TestReview_CommitsFlashcardAndReviewTogether(t *testing.T) {
	resetDB(t)
	ctx := testContext(t)
	flashcard := createFlashcard(t, createDeck(t))

	repo := flashcardrepository.NewFlashcardRepository(testDB)
	svc := service.NewFlashcardService(repo, reviewrepository.NewReviewRepository(testDB), gormx.NewTxManager(testDB))

	require.NoError(t, svc.Review(ctx, flashcard.ID.String(), "good"))

	updated, err := repo.FindByID(ctx, flashcard.ID.String())
	require.NoError(t, err)
	assert.Equal(t, domain.StateLearning, updated.State)
	assert.True(t, updated.LastReviewAt.Valid)

	var review struct {
		Response      string
		PreviousState string
		NewState      string
	}
	require.NoError(t, testDB.Table("reviews").Where("flashcard_id = ?", flashcard.ID.String()).Take(&review).Error)
	assert.Equal(t, "good", review.Response)
	assert.Equal(t, "new", review.PreviousState)
	assert.Equal(t, "learning", review.NewState)
}

func TestReview_RollsBackFlashcardWhenReviewInsertFails(t *testing.T) {
	resetDB(t)
	ctx := testContext(t)
	flashcard := createFlashcard(t, createDeck(t))

	repo := flashcardrepository.NewFlashcardRepository(testDB)
	svc := service.NewFlashcardService(repo, failingReviewRepository{}, gormx.NewTxManager(testDB))

	err := svc.Review(ctx, flashcard.ID.String(), "good")
	require.ErrorIs(t, err, errReviewInsert)

	unchanged, err := repo.FindByID(ctx, flashcard.ID.String())
	require.NoError(t, err)
	assert.Equal(t, domain.StateNew, unchanged.State)
	assert.False(t, unchanged.LastReviewAt.Valid)
	assert.False(t, unchanged.NextReviewAt.Valid)
	assert.Zero(t, unchanged.Interval)
	assert.Zero(t, countReviews(t, flashcard.ID.String()))
}
