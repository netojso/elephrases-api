package portrepository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type DeckRepository interface {
	FindAll(ctx context.Context) ([]*domain.Deck, error)
	FindByID(ctx context.Context, id string) (*domain.Deck, error)
	Save(ctx context.Context, Deck *domain.Deck) error
	Update(ctx context.Context, Deck *domain.Deck) error
	Delete(ctx context.Context, id string) error
}
