package domain

import (
	"time"

	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type Deck struct {
	ID          pkg.UUID                `json:"id"`
	Name        string                  `json:"name"`
	Description nullable.NullableString `json:"description"`
	Category    string                  `json:"category"`
	Visibility  string                  `json:"visibility"`
	CreatedAt   time.Time               `json:"created_at"`
	Flashcards  []Flashcard             `json:"flashcards"`
}

func NewDeck(name string, description string, category string, visibility string) *Deck {
	return &Deck{
		ID:          pkg.NewUUID(),
		Name:        name,
		Description: nullable.NewNullableString(description),
		Category:    category,
		Visibility:  visibility,
		CreatedAt:   time.Now(),
	}
}
