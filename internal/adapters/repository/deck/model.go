package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	repository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type Deck struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	Name        string                 `gorm:"type:varchar(255)" json:"name"`
	Description sql.NullString         `gorm:"type:text" json:"description"`
	Category    string                 `gorm:"type:varchar(255)" json:"category"`
	Visibility  string                 `gorm:"type:varchar(255)" json:"visibility"`
	CreatedAt   time.Time              `gorm:"type:timestamp" json:"created_at"`
	Flashcards  []repository.Flashcard `gorm:"foreignKey:DeckID" json:"flashcards"`
}

func (d *Deck) ToDomain() *domain.Deck {
	id, _ := pkg.ParseUUID(d.ID.String())

	flashbacks := make([]domain.Flashcard, len(d.Flashcards))
	for i, f := range d.Flashcards {
		flashbacks[i] = *f.ToDomain()
	}

	return &domain.Deck{
		ID:          id,
		Name:        d.Name,
		Description: nullable.NewNullableString(d.Description.String),
		Category:    d.Category,
		Visibility:  d.Visibility,
		CreatedAt:   d.CreatedAt,
		Flashcards:  flashbacks,
	}
}

func domainToModel(deck *domain.Deck) *Deck {
	return &Deck{
		ID:          deck.ID.Value(),
		Name:        deck.Name,
		Description: deck.Description.NullString,
		Category:    deck.Category,
		Visibility:  deck.Visibility,
		CreatedAt:   deck.CreatedAt,
	}
}

func modelToDomainSlice(decks []*Deck) []*domain.Deck {
	var result []*domain.Deck
	for _, deck := range decks {
		result = append(result, deck.ToDomain())
	}
	return result
}
