//go:build integration

package integration

import (
	"testing"

	deckrepository "github.com/netojso/elephrases-api/internal/adapters/repository/deck"
	flashcardrepository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/stretchr/testify/require"
)

func createDeck(t *testing.T) *domain.Deck {
	t.Helper()
	deck := domain.NewDeck("deck", "", "category", "private")
	require.NoError(t, deckrepository.NewDeckRepository(testDB).Save(testContext(t), deck))
	return deck
}

func createFlashcard(t *testing.T, deck *domain.Deck, mutate ...func(*domain.Flashcard)) *domain.Flashcard {
	t.Helper()
	flashcard := domain.NewFlashcard(deck.ID, "front", "back", "")
	for _, m := range mutate {
		m(flashcard)
	}
	require.NoError(t, flashcardrepository.NewFlashcardRepository(testDB).Save(testContext(t), flashcard))
	return flashcard
}

func countReviews(t *testing.T, flashcardID string) int64 {
	t.Helper()
	var n int64
	require.NoError(t, testDB.Table("reviews").Where("flashcard_id = ?", flashcardID).Count(&n).Error)
	return n
}
