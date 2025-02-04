package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal"
)

type CardState string

const (
	New      CardState = "new"
	Learning CardState = "learning"
	Review   CardState = "review"
	Lapsed   CardState = "lapsed"
)

type Flashcard struct {
	ID           uuid.UUID             `json:"id"`
	DeckID       uuid.UUID             `json:"deck_id" validate:"required"`
	Front        string                `json:"front" validate:"required"`
	Back         string                `json:"back" validate:"required"`
	CreatedAt    time.Time             `json:"created_at" validate:"required"`
	LastReviewAt internal.NullableTime `json:"last_review_at"`
	NextReviewAt internal.NullableTime `json:"next_review_at"`
	State        CardState             `json:"state"`
	EaseFactor   float64               `json:"ease_factor"`
	Interval     time.Duration         `json:"interval"`
}

type CreateFlashcardRequest struct {
	DeckID string `json:"deck_id" binding:"required"`
	Front  string `json:"front" binding:"required"`
	Back   string `json:"back" binding:"required"`
}

func (f *Flashcard) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func NewFlashcard(
	params CreateFlashcardRequest,
) (*Flashcard, error) {
	flashcard := &Flashcard{
		ID:        uuid.New(),
		DeckID:    uuid.MustParse(params.DeckID),
		Front:     params.Front,
		Back:      params.Back,
		CreatedAt: time.Now(),
		State:     New,
	}

	if err := flashcard.Validate(); err != nil {
		return nil, err
	}

	return flashcard, nil
}

func (f *Flashcard) ReviewFlashcard(response string, settings *Settings) {
	f.LastReviewAt = internal.NewNullableTime(time.Now())

	if f.State == New {
		f.State = Learning
	}

	if settings == nil {
		settings = &Settings{
			LearningSteps:      []time.Duration{time.Minute, 10 * time.Minute},
			GraduatingInterval: 24 * time.Hour,
			EasyInterval:       4 * 24 * time.Hour,
		}
	}
	switch response {
	case "again":
		f.Interval = settings.LearningSteps[0]
		f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
	case "good":
		if f.State == New {
			f.Interval = settings.LearningSteps[0]
			f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
		}

		if f.State == Learning {
			if f.Interval == settings.LearningSteps[len(settings.LearningSteps)-1] {
				f.State = Review
				f.Interval = settings.GraduatingInterval
				f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
			} else {
				find_index := -1
				for i, v := range settings.LearningSteps {
					if v == f.Interval {
						find_index = i
						break
					}
				}
				f.Interval = settings.LearningSteps[find_index+1]
				f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
			}
		}
	case "hard":
		fmt.Println(f.State)
		if f.State == Learning {
			fmt.Println(settings.LearningSteps)
			if f.Interval == settings.LearningSteps[0] {
				if len(settings.LearningSteps) > 1 {
					average := (settings.LearningSteps[0] + settings.LearningSteps[1]) / 2
					fmt.Println(average)
					f.Interval = average
					f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
				} else {
					hard_delay := time.Duration(1.5 * float64(settings.LearningSteps[0]))
					if hard_delay > settings.LearningSteps[0]+24*time.Hour {
						hard_delay = settings.LearningSteps[0] + 24*time.Hour
					}

					f.Interval = hard_delay
					f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
				}
			} else {
				f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
			}
		}

	case "easy":
		f.State = Review
		f.Interval = settings.EasyInterval
		f.NextReviewAt = internal.NewNullableTime(time.Now().Add(f.Interval))
	}

}

type Options struct {
	Where map[string]interface{}
}

type FlashcardRepository interface {
	FindAll(options *Options) ([]Flashcard, error)
	FindByID(id string) (Flashcard, error)
	Save(flashcard Flashcard) error
	Update(flashcard Flashcard) error
	Delete(id string) error
}

type FlashcardUsecase interface {
	GetDueFlashcards() ([]Flashcard, error)
	GetAll() ([]Flashcard, error)
	GetByID(id string) (Flashcard, error)
	Create(flashcard Flashcard) error
	Update(flashcard Flashcard) error
	Review(id string, response string) error
	Delete(id string) error
}
