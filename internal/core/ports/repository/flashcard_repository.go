package portrepository

import "github.com/netojso/elephrases-api/internal/core/domain"

type Options struct {
	Where map[string]interface{}
}

type FlashcardRepository interface {
	FindByDeckID(deckID string) ([]*domain.Flashcard, error)
	FindAll(options *Options) ([]*domain.Flashcard, error)
	FindByID(id string) (*domain.Flashcard, error)
	Save(flashcard *domain.Flashcard) error
	Update(flashcard *domain.Flashcard) error
	Delete(id string) error
}
