//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	deckrepository "github.com/netojso/elephrases-api/internal/adapters/repository/deck"
	flashcardrepository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
	userrepository "github.com/netojso/elephrases-api/internal/adapters/repository/users"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlashcardConstraints(t *testing.T) {
	resetDB(t)
	deck := createDeck(t)
	repo := flashcardrepository.NewFlashcardRepository(testDB)

	tests := []struct {
		name       string
		mutate     func(*domain.Flashcard)
		constraint string
	}{
		{"state inválido", func(f *domain.Flashcard) { f.State = "archived" }, "chk_flashcards_state"},
		{"ease_factor abaixo de 1.3", func(f *domain.Flashcard) { f.EaseFactor = 1.0 }, "chk_flashcards_ease"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flashcard := domain.NewFlashcard(deck.ID, "front", "back", "")
			tt.mutate(flashcard)

			err := repo.Save(testContext(t), flashcard)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.constraint)
		})
	}

	t.Run("state lapsed é aceito", func(t *testing.T) {
		flashcard := domain.NewFlashcard(deck.ID, "front", "back", "")
		flashcard.State = "lapsed"
		assert.NoError(t, repo.Save(testContext(t), flashcard))
	})
}

func TestDeleteDeck_CascadesToFlashcardsAndReviews(t *testing.T) {
	resetDB(t)
	ctx := testContext(t)
	deck := createDeck(t)
	flashcard := createFlashcard(t, deck)
	require.NoError(t, testDB.Exec(
		"INSERT INTO reviews (flashcard_id, response, previous_state, new_state, interval_ns, ease_factor) VALUES (?, 'good', 'new', 'learning', 0, 2.5)",
		flashcard.ID.String(),
	).Error)

	require.NoError(t, deckrepository.NewDeckRepository(testDB).Delete(ctx, deck.ID.String()))

	_, err := flashcardrepository.NewFlashcardRepository(testDB).FindByID(ctx, flashcard.ID.String())
	assert.ErrorIs(t, err, pkg.ErrNotFound)
	assert.Zero(t, countReviews(t, flashcard.ID.String()))
}

func TestRepositories_MapNotFound(t *testing.T) {
	resetDB(t)
	ctx := testContext(t)
	missing := pkg.NewUUID().String()

	flashcards := flashcardrepository.NewFlashcardRepository(testDB)
	decks := deckrepository.NewDeckRepository(testDB)
	users := userrepository.NewUserRepository(testDB)

	checks := map[string]func() error{
		"flashcard FindByID":          func() error { _, err := flashcards.FindByID(ctx, missing); return err },
		"flashcard FindByIDForUpdate": func() error { _, err := flashcards.FindByIDForUpdate(ctx, missing); return err },
		"flashcard Delete":            func() error { return flashcards.Delete(ctx, missing) },
		"deck FindByID":               func() error { _, err := decks.FindByID(ctx, missing); return err },
		"deck Delete":                 func() error { return decks.Delete(ctx, missing) },
		"user GetByID":                func() error { _, err := users.GetByID(ctx, missing); return err },
		"user GetByEmail":             func() error { _, err := users.GetByEmail(ctx, "missing@x.com"); return err },
		"user DeleteUser":             func() error { return users.DeleteUser(ctx, missing) },
	}

	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, check(), pkg.ErrNotFound)
		})
	}
}

func TestFindByIDForUpdate_BlocksConcurrentLock(t *testing.T) {
	resetDB(t)
	flashcard := createFlashcard(t, createDeck(t))
	repo := flashcardrepository.NewFlashcardRepository(testDB)
	txm := gormx.NewTxManager(testDB)

	locked := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)

	// Transação 1 segura o lock até receber release.
	go func() {
		done <- txm.WithinTransaction(context.Background(), func(ctx context.Context) error {
			if _, err := repo.FindByIDForUpdate(ctx, flashcard.ID.String()); err != nil {
				return err
			}
			close(locked)
			<-release
			return nil
		})
	}()
	<-locked

	// Transação 2 não consegue o lock enquanto a 1 estiver aberta.
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	err := txm.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err := repo.FindByIDForUpdate(ctx, flashcard.ID.String())
		return err
	})
	assert.Error(t, err, "o segundo SELECT ... FOR UPDATE deveria esperar o lock")

	close(release)
	require.NoError(t, <-done)

	// Com o lock liberado, a leitura bloqueante volta a funcionar.
	require.NoError(t, txm.WithinTransaction(testContext(t), func(ctx context.Context) error {
		_, err := repo.FindByIDForUpdate(ctx, flashcard.ID.String())
		return err
	}))
}
