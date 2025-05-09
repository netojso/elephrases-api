package domain

import (
	"time"

	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type Stats struct {
	NewCards      int `json:"new_cards"`
	LearningCards int `json:"learning_cards"`
	ReviewCards   int `json:"review_cards"`
	TotalCards    int `json:"total_cards"`
}
type Deck struct {
	ID          pkg.UUID                `json:"id"`
	Name        string                  `json:"name"`
	Description nullable.NullableString `json:"description"`
	Category    string                  `json:"category"`
	Visibility  string                  `json:"visibility"`
	CreatedAt   time.Time               `json:"created_at"`
	Flashcards  []Flashcard             `json:"-"`
	Stats       Stats                   `json:"stats"`
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
