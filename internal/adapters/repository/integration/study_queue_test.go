//go:build integration

package integration

import (
	"testing"
	"time"

	flashcardrepository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/pkg/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fronts(cards []*domain.Flashcard) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.Front
	}
	return out
}

func TestFindStudyQueue(t *testing.T) {
	resetDB(t)
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)

	deck := createDeck(t)
	other := createDeck(t)

	newCard := func(front string, createdAt time.Time) func(*domain.Flashcard) {
		return func(f *domain.Flashcard) {
			f.Front = front
			f.CreatedAt = createdAt
		}
	}
	reviewed := func(front string, lastReview, nextReview time.Time) func(*domain.Flashcard) {
		return func(f *domain.Flashcard) {
			f.Front = front
			f.State = domain.StateReview
			f.LastReviewAt = nullable.NewNullableTime(lastReview)
			f.NextReviewAt = nullable.NewNullableTime(nextReview)
		}
	}

	// Novos, criados em ordem new-1 < new-2 < new-3.
	createFlashcard(t, deck, newCard("new-3", now.Add(-1*time.Hour)))
	createFlashcard(t, deck, newCard("new-1", now.Add(-3*time.Hour)))
	createFlashcard(t, deck, newCard("new-2", now.Add(-2*time.Hour)))
	// Vencidos, revisados antes de hoje.
	createFlashcard(t, deck, reviewed("due-late", now.Add(-48*time.Hour), now.Add(-1*time.Hour)))
	createFlashcard(t, deck, reviewed("due-early", now.Add(-72*time.Hour), now.Add(-5*time.Hour)))
	// Revisado hoje e ainda não vencido: conta para o limite diário.
	createFlashcard(t, deck, reviewed("today", now.Add(-2*time.Hour), now.Add(24*time.Hour)))
	// Revisado ontem e não vencido: fora da fila.
	createFlashcard(t, deck, reviewed("future", now.Add(-20*time.Hour), now.Add(48*time.Hour)))
	// Outro deck: nunca aparece.
	createFlashcard(t, other, newCard("other-new", now.Add(-4*time.Hour)))
	createFlashcard(t, other, reviewed("other-due", now.Add(-48*time.Hour), now.Add(-1*time.Hour)))

	repo := flashcardrepository.NewFlashcardRepository(testDB)

	tests := []struct {
		name       string
		dailyLimit int
		wantNew    []string
	}{
		{name: "limite restante = 3 - 1 revisado hoje", dailyLimit: 3, wantNew: []string{"new-1", "new-2"}},
		{name: "limite maior que os novos disponíveis", dailyLimit: 20, wantNew: []string{"new-1", "new-2", "new-3"}},
		{name: "limite já atingido não devolve novos", dailyLimit: 1, wantNew: []string{}},
		{name: "limite excedido não vira 'sem limite'", dailyLimit: 0, wantNew: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			study, err := repo.FindStudyQueue(testContext(t), deck.ID.String(), tt.dailyLimit, now)
			require.NoError(t, err)

			assert.Equal(t, tt.wantNew, fronts(study.NewCards))
			assert.Equal(t, []string{"due-early", "due-late"}, fronts(study.DueCards))
		})
	}
}

func TestFindStudyQueue_EmptyDeckReturnsEmptySlices(t *testing.T) {
	resetDB(t)
	deck := createDeck(t)

	study, err := flashcardrepository.NewFlashcardRepository(testDB).FindStudyQueue(testContext(t), deck.ID.String(), 20, time.Now())
	require.NoError(t, err)
	assert.NotNil(t, study.NewCards)
	assert.NotNil(t, study.DueCards)
	assert.Empty(t, study.NewCards)
	assert.Empty(t, study.DueCards)
}
