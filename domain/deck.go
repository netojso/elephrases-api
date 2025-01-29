package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal"
)

type Settings struct {
	LearningSteps      []time.Duration `json:"learning_steps"`
	GraduatingInterval time.Duration   `json:"graduating_interval"`
	EasyInterval       time.Duration   `json:"easy_interval"`
}

type Deck struct {
	ID          uuid.UUID               `json:"id" validate:"required"`
	Name        string                  `json:"name" validate:"required"`
	Description internal.NullableString `json:"description" validate:"required"`
	Category    string                  `json:"category" validate:"required"`
	Visibility  string                  `json:"visibility" validate:"required"`
	CreatedAt   time.Time               `json:"created_at" validate:"required"`
	Flashcards  []Flashcard             `json:"flashcards"`
	// Settings    *Seetings               `json:"settings"`
}

type CreateDeckRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Description internal.NullableString `json:"description" binding:"required"`
	Category    string                  `json:"category" binding:"required"`
	Visibility  string                  `json:"visibility" binding:"required"`
}

func (f *Deck) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func NewDeck(
	params CreateDeckRequest,
) (*Deck, error) {
	deck := &Deck{
		ID:          uuid.New(),
		Name:        params.Name,
		Description: params.Description,
		Category:    params.Category,
		Visibility:  params.Visibility,
		CreatedAt:   time.Now(),
	}

	if err := deck.Validate(); err != nil {
		return nil, err
	}

	return deck, nil
}

type DeckRepository interface {
	FindAll() ([]Deck, error)
	FindByID(id string) (Deck, error)
	Save(Deck Deck) error
	Update(Deck Deck) error
	Delete(id string) error
}

type DeckUsecase interface {
	GetAll() ([]Deck, error)
	GetByID(id string) (Deck, error)
	Create(Deck Deck) error
	Update(Deck Deck) error
	Delete(id string) error
}
