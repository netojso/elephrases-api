package portrepository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type Options struct {
	Where map[string]interface{}
	Limit int
}

type FlashcardRepository interface {
	FindByDeckID(ctx context.Context, deckID string) ([]*domain.Flashcard, error)
	FindAll(ctx context.Context, options *Options) ([]*domain.Flashcard, error)
	FindByID(ctx context.Context, id string) (*domain.Flashcard, error)
	Save(ctx context.Context, flashcard *domain.Flashcard) error
	Update(ctx context.Context, flashcard *domain.Flashcard) error
	Delete(ctx context.Context, id string) error
}
