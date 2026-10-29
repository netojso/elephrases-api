package portservice

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type DeckService interface {
	GetAll(ctx context.Context) ([]*domain.Deck, error)
	GetByID(ctx context.Context, id string) (*domain.Deck, error)
	Create(ctx context.Context, Deck *domain.Deck) error
	Update(ctx context.Context, Deck *domain.Deck) error
	Delete(ctx context.Context, id string) error
}
