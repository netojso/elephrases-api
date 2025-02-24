package portservice

import "github.com/netojso/elephrases-api/internal/core/domain"

type FlashcardService interface {
	GetDueFlashcards() ([]*domain.Flashcard, error)
	GetAll() ([]*domain.Flashcard, error)
	GetByID(id string) (*domain.Flashcard, error)
	GetByDeckID(deckID string) ([]*domain.Flashcard, error)
	Create(flashcard *domain.Flashcard) error
	Update(flashcard *domain.Flashcard) error
	Review(id string, response string) error
	Delete(id string) error
}
