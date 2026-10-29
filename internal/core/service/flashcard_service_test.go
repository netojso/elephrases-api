package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/internal/core/service"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var fixedNow = time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)

type flashcardDeps struct {
	repo    *flashcardRepoMock
	reviews *reviewRepoMock
	tx      *txManagerMock
}

func newFlashcardService(t *testing.T) (*flashcardDeps, portservice.FlashcardService) {
	t.Helper()
	deps := &flashcardDeps{repo: &flashcardRepoMock{}, reviews: &reviewRepoMock{}, tx: &txManagerMock{}}
	t.Cleanup(func() {
		deps.repo.AssertExpectations(t)
		deps.reviews.AssertExpectations(t)
		deps.tx.AssertExpectations(t)
	})
	svc := service.NewFlashcardService(deps.repo, deps.reviews, deps.tx, service.WithClock(func() time.Time { return fixedNow }))
	return deps, svc
}

func TestFlashcardService_Study_UsesDailyLimitAndClock(t *testing.T) {
	deps, svc := newFlashcardService(t)
	want := &domain.Study{NewCards: []*domain.Flashcard{}, DueCards: []*domain.Flashcard{}}
	deps.repo.On("FindStudyQueue", mock.Anything, "deck-1", 20, fixedNow).Return(want, nil)

	got, err := svc.Study(context.Background(), "deck-1")

	require.NoError(t, err)
	assert.Same(t, want, got)
}

func TestFlashcardService_Review_UpdatesCardAndRecordsReviewInTransaction(t *testing.T) {
	deps, svc := newFlashcardService(t)
	card := domain.NewFlashcard(pkg.NewUUID(), "front", "back", "")

	deps.tx.On("WithinTransaction", mock.Anything).Once()
	deps.repo.On("FindByIDForUpdate", mock.MatchedBy(inTx), card.ID.String()).Return(card, nil)
	deps.repo.On("Update", mock.MatchedBy(inTx), mock.MatchedBy(func(f *domain.Flashcard) bool {
		return f.State == domain.StateLearning && f.LastReviewAt.Time.Equal(fixedNow)
	})).Return(nil)
	deps.reviews.On("Save", mock.MatchedBy(inTx), mock.MatchedBy(func(r *domain.Review) bool {
		return r.FlashcardID == card.ID &&
			r.Response == "good" &&
			r.PreviousState == domain.StateNew &&
			r.NewState == domain.StateLearning &&
			r.Interval == 10*time.Minute &&
			r.EaseFactor == domain.DefaultEaseFactor &&
			r.ReviewedAt.Equal(fixedNow)
	})).Return(nil)

	require.NoError(t, svc.Review(context.Background(), card.ID.String(), "good"))
}

func TestFlashcardService_Review_NotFound(t *testing.T) {
	deps, svc := newFlashcardService(t)
	deps.tx.On("WithinTransaction", mock.Anything).Once()
	deps.repo.On("FindByIDForUpdate", mock.Anything, "missing").Return(nil, pkg.ErrNotFound)

	err := svc.Review(context.Background(), "missing", "good")

	assert.ErrorIs(t, err, pkg.ErrNotFound)
	deps.repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	deps.reviews.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestFlashcardService_Review_UpdateFailureSkipsReviewInsert(t *testing.T) {
	deps, svc := newFlashcardService(t)
	card := domain.NewFlashcard(pkg.NewUUID(), "front", "back", "")
	boom := errors.New("update failed")

	deps.tx.On("WithinTransaction", mock.Anything).Once()
	deps.repo.On("FindByIDForUpdate", mock.Anything, card.ID.String()).Return(card, nil)
	deps.repo.On("Update", mock.Anything, mock.Anything).Return(boom)

	err := svc.Review(context.Background(), card.ID.String(), "good")

	assert.ErrorIs(t, err, boom)
	deps.reviews.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestFlashcardService_Review_PropagatesReviewInsertFailure(t *testing.T) {
	deps, svc := newFlashcardService(t)
	card := domain.NewFlashcard(pkg.NewUUID(), "front", "back", "")
	boom := errors.New("insert failed")

	deps.tx.On("WithinTransaction", mock.Anything).Once()
	deps.repo.On("FindByIDForUpdate", mock.Anything, card.ID.String()).Return(card, nil)
	deps.repo.On("Update", mock.Anything, mock.Anything).Return(nil)
	deps.reviews.On("Save", mock.Anything, mock.Anything).Return(boom)

	// O erro precisa chegar ao TxManager para que ele faça o rollback.
	assert.ErrorIs(t, svc.Review(context.Background(), card.ID.String(), "good"), boom)
}

func TestFlashcardService_Create_PrefixesMediaWithS3URL(t *testing.T) {
	viper.Set("AWS_S3_URL", "https://cdn.example.com/")
	t.Cleanup(func() { viper.Set("AWS_S3_URL", "") })

	tests := []struct {
		name      string
		media     string
		wantValid bool
		wantURL   string
	}{
		{name: "com mídia", media: "cat.png", wantValid: true, wantURL: "https://cdn.example.com/cat.png"},
		{name: "sem mídia", media: "", wantValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, svc := newFlashcardService(t)
			card := domain.NewFlashcard(pkg.NewUUID(), "front", "back", tt.media)
			deps.repo.On("Save", mock.Anything, card).Return(nil)

			require.NoError(t, svc.Create(context.Background(), card))

			assert.Equal(t, tt.wantValid, card.MediaUrl.Valid)
			if tt.wantValid {
				assert.Equal(t, tt.wantURL, card.MediaUrl.String)
			}
		})
	}
}
