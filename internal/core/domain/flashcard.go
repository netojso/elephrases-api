package domain

import (
	"fmt"
	"time"

	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type CardState string

const (
	New      CardState = "new"
	Learning CardState = "learning"
	Review   CardState = "review"
	Lapsed   CardState = "lapsed"
)

type Flashcard struct {
	ID           pkg.UUID              `json:"id"`
	DeckID       pkg.UUID              `json:"deck_id"`
	Front        string                `json:"front"`
	Back         string                `json:"back"`
	CreatedAt    time.Time             `json:"created_at"`
	LastReviewAt nullable.NullableTime `json:"last_review_at"`
	NextReviewAt nullable.NullableTime `json:"next_review_at"`
	State        CardState             `json:"state"`
	EaseFactor   float64               `json:"ease_factor"`
	Interval     time.Duration         `json:"interval"`
}

type Settings struct {
	LearningSteps      []time.Duration
	GraduatingInterval time.Duration
	EasyInterval       time.Duration
}

func NewFlashcard(deckID pkg.UUID, front string, back string) *Flashcard {
	return &Flashcard{
		ID:        pkg.NewUUID(),
		DeckID:    deckID,
		Front:     front,
		Back:      back,
		CreatedAt: time.Now(),
		State:     New,
	}
}

func (f *Flashcard) ReviewFlashcard(response string, settings *Settings) {
	f.LastReviewAt = nullable.NewNullableTime(time.Now())

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
		f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
	case "good":
		if f.State == New {
			f.Interval = settings.LearningSteps[0]
			f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
		}

		if f.State == Learning {
			if f.Interval == settings.LearningSteps[len(settings.LearningSteps)-1] {
				f.State = Review
				f.Interval = settings.GraduatingInterval
				f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
			} else {
				find_index := -1
				for i, v := range settings.LearningSteps {
					if v == f.Interval {
						find_index = i
						break
					}
				}
				f.Interval = settings.LearningSteps[find_index+1]
				f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
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
					f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
				} else {
					hard_delay := time.Duration(1.5 * float64(settings.LearningSteps[0]))
					if hard_delay > settings.LearningSteps[0]+24*time.Hour {
						hard_delay = settings.LearningSteps[0] + 24*time.Hour
					}

					f.Interval = hard_delay
					f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
				}
			} else {
				f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
			}
		}

	case "easy":
		f.State = Review
		f.Interval = settings.EasyInterval
		f.NextReviewAt = nullable.NewNullableTime(time.Now().Add(f.Interval))
	}

}
