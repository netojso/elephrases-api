package domain

import (
	"log"
	"time"

	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type CardState string

const (
	StateNew      CardState = "new"
	StateLearning CardState = "learning"
	StateReview   CardState = "review"
)

type Flashcard struct {
	ID           pkg.UUID                `json:"id"`
	DeckID       pkg.UUID                `json:"deck_id"`
	MediaUrl     nullable.NullableString `json:"media_url"`
	Front        string                  `json:"front"`
	Back         string                  `json:"back"`
	CreatedAt    time.Time               `json:"created_at"`
	LastReviewAt nullable.NullableTime   `json:"last_review_at"`
	NextReviewAt nullable.NullableTime   `json:"next_review_at"`
	State        CardState               `json:"state"`
	EaseFactor   float64                 `json:"ease_factor"`
	Interval     time.Duration           `json:"interval"`
}

type Settings struct {
	LearningSteps      []time.Duration
	GraduatingInterval time.Duration
	EasyInterval       time.Duration
}

func NewFlashcard(deckID pkg.UUID, front string, back string, media string) *Flashcard {
	return &Flashcard{
		ID:        pkg.NewUUID(),
		DeckID:    deckID,
		Front:     front,
		Back:      back,
		CreatedAt: time.Now(),
		State:     StateNew,
		MediaUrl:  nullable.NewNullableString(media),
	}
}

func (f *Flashcard) ReviewFlashcard(response string, settings *Settings) {
	// Set last review time
	now := time.Now()
	f.LastReviewAt = nullable.NewNullableTime(now)

	log.Println(response)

	// Default settings if nil
	if settings == nil {
		settings = &Settings{
			LearningSteps:      []time.Duration{time.Minute, 10 * time.Minute},
			GraduatingInterval: 24 * time.Hour,
			EasyInterval:       4 * 24 * time.Hour,
		}
	}

	// Transition New to Learning and initialize interval
	if f.State == StateNew {
		f.State = StateLearning
		f.Interval = settings.LearningSteps[0] // Start at 1 minute
	}

	// Handle response
	switch response {
	case "again":
		f.Interval = settings.LearningSteps[0]
		f.NextReviewAt = nullable.NewNullableTime(now.Add(f.Interval))

	case "good":
		log.Println(f)
		if f.State == StateLearning {
			currentStep := f.findLearningStepIndex(settings.LearningSteps)
			log.Println(currentStep)
			if currentStep == len(settings.LearningSteps)-1 {
				f.State = StateReview
				f.Interval = settings.GraduatingInterval
				f.NextReviewAt = nullable.NewNullableTime(now.Add(f.Interval))
			} else {
				f.Interval = settings.LearningSteps[currentStep+1]
				f.NextReviewAt = nullable.NewNullableTime(now.Add(f.Interval))
			}
		}

	case "hard":
		if f.State == StateLearning {
			if f.Interval == settings.LearningSteps[0] && len(settings.LearningSteps) > 1 {
				f.Interval = (settings.LearningSteps[0] + settings.LearningSteps[1]) / 2
			} else if len(settings.LearningSteps) == 1 {
				f.Interval = time.Duration(1.5 * float64(settings.LearningSteps[0]))
				if f.Interval > settings.LearningSteps[0]+24*time.Hour {
					f.Interval = settings.LearningSteps[0] + 24*time.Hour
				}
			}
			f.NextReviewAt = nullable.NewNullableTime(now.Add(f.Interval))
		}

	case "easy":
		f.State = StateReview
		f.Interval = settings.EasyInterval
		f.NextReviewAt = nullable.NewNullableTime(now.Add(f.Interval))

	default:
		f.NextReviewAt = nullable.NewNullableTime(now.Add(f.Interval))
	}
}

func (f *Flashcard) findLearningStepIndex(steps []time.Duration) int {
	for i, step := range steps {
		if f.Interval == step {
			return i
		}
	}
	return 0 // Default to first step if not found
}
