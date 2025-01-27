package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal"
)

type Flashcard struct {
	ID             uuid.UUID               `json:"id"`
	Front          string                  `json:"front" binding:"required" validate:"required"`
	Back           string                  `json:"back" binding:"required" validate:"required"`
	CreatedAt      string                  `json:"created_at" validate:"required"`
	LastReviewedAt string                  `json:"last_reviewed_at"`
	Status         internal.NullableString `json:"status"`
}

type CreateFlashcardRequest struct {
	Front string `json:"front" validate:"required"`
	Back  string `json:"back" validate:"required"`
}

func (f *Flashcard) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func NewFlashcard(
	params CreateFlashcardRequest,
) (*Flashcard, error) {
	flashcard := &Flashcard{
		ID:             uuid.New(),
		Front:          params.Front,
		Back:           params.Back,
		CreatedAt:      time.Now().Format(time.RFC3339),
		LastReviewedAt: time.Now().Format(time.RFC3339),
	}

	if err := flashcard.Validate(); err != nil {
		return nil, err
	}

	return flashcard, nil
}

type FlashcardRepository interface {
	FindAll() ([]Flashcard, error)
	FindByID(id string) (Flashcard, error)
	Save(flashcard Flashcard) error
	Update(flashcard Flashcard) error
	Delete(id string) error
}

type FlashcardUsecase interface {
	GetAll() ([]Flashcard, error)
	GetByID(id string) (Flashcard, error)
	Create(flashcard Flashcard) error
	Update(flashcard Flashcard) error
	Delete(id string) error
}
