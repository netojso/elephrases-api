package domain

import (
	"time"

	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type Deck struct {
	ID          pkg.UUID
	Name        string
	Description nullable.NullableString
	Category    string
	Visibility  string
	CreatedAt   time.Time
	Flashcards  []Flashcard
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
