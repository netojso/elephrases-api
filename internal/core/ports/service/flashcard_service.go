package portservice

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type FlashcardService interface {
	Study(ctx context.Context, deckID string) (*domain.Study, error)
	GetDueFlashcards(ctx context.Context) ([]*domain.Flashcard, error)
	GetAll(ctx context.Context) ([]*domain.Flashcard, error)
	GetByID(ctx context.Context, id string) (*domain.Flashcard, error)
	GetByDeckID(ctx context.Context, deckID string) ([]*domain.Flashcard, error)
	Create(ctx context.Context, flashcard *domain.Flashcard) error
	Update(ctx context.Context, flashcard *domain.Flashcard) error
	Review(ctx context.Context, id string, response string) error
	Delete(ctx context.Context, id string) error
}
