package domain

import (
	"time"

	"github.com/netojso/elephrases-api/pkg"
)

// Review registra uma resposta do usuário a um cartão e o estado resultante.
type Review struct {
	ID            pkg.UUID
	FlashcardID   pkg.UUID
	Response      string
	PreviousState CardState
	NewState      CardState
	Interval      time.Duration
	EaseFactor    float64
	ReviewedAt    time.Time
}

func NewReview(flashcard *Flashcard, response string, previousState CardState, reviewedAt time.Time) *Review {
	return &Review{
		ID:            pkg.NewUUID(),
		FlashcardID:   flashcard.ID,
		Response:      response,
		PreviousState: previousState,
		NewState:      flashcard.State,
		Interval:      flashcard.Interval,
		EaseFactor:    flashcard.EaseFactor,
		ReviewedAt:    reviewedAt,
	}
}
