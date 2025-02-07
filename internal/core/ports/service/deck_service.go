package portservice

import "github.com/netojso/elephrases-api/internal/core/domain"

type DeckService interface {
	GetAll() ([]*domain.Deck, error)
	GetByID(id string) (*domain.Deck, error)
	Create(Deck *domain.Deck) error
	Update(Deck *domain.Deck) error
	Delete(id string) error
}
