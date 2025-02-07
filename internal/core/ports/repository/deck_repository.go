package portrepository

import "github.com/netojso/elephrases-api/internal/core/domain"

type DeckRepository interface {
	FindAll() ([]*domain.Deck, error)
	FindByID(id string) (*domain.Deck, error)
	Save(Deck *domain.Deck) error
	Update(Deck *domain.Deck) error
	Delete(id string) error
}
