package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal/core/domain"
	"github.com/netojso/elephrases-api/pkg"
	"github.com/netojso/elephrases-api/pkg/nullable"
)

type Deck struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	Name string    `gorm:"type:varchar(100);not null"`
}

type Flashcard struct {
	ID           uuid.UUID     `gorm:"type:uuid;primary_key"`
	Front        string        `gorm:"type:text"`
	Back         string        `gorm:"type:text"`
	CreatedAt    time.Time     `gorm:"type:timestamp"`
	LastReviewAt sql.NullTime  `gorm:"type:timestamp"`
	NextReviewAt sql.NullTime  `gorm:"type:timestamp"`
	State        string        `gorm:"type:varchar(20)"`
	EaseFactor   float64       `gorm:"type:float"`
	Interval     time.Duration `gorm:"type:interval"`
	DeckID       uuid.UUID     `gorm:"foreignKey:DeckID;references:ID"`
}

func (f *Flashcard) TableName() string {
	return "flashcards"
}

func (f *Flashcard) ToDomain() *domain.Flashcard {
	id, err := pkg.ParseUUID(f.ID.String())

	if err != nil {
		return nil
	}
	deckID, err := pkg.ParseUUID(f.DeckID.String())
	if err != nil {
		return nil
	}

	return &domain.Flashcard{
		ID:           id,
		DeckID:       deckID,
		Front:        f.Front,
		Back:         f.Back,
		CreatedAt:    f.CreatedAt,
		LastReviewAt: nullable.NewNullableTime(f.LastReviewAt.Time),
		NextReviewAt: nullable.NewNullableTime(f.NextReviewAt.Time),
		State:        domain.CardState(f.State),
		EaseFactor:   f.EaseFactor,
		Interval:     f.Interval,
	}
}

func domainToModel(flashcard *domain.Flashcard) *Flashcard {
	return &Flashcard{
		ID:           flashcard.ID.Value(),
		DeckID:       flashcard.DeckID.Value(),
		Front:        flashcard.Front,
		Back:         flashcard.Back,
		CreatedAt:    flashcard.CreatedAt,
		LastReviewAt: flashcard.LastReviewAt.NullTime,
		NextReviewAt: flashcard.NextReviewAt.NullTime,
		State:        string(flashcard.State),
		EaseFactor:   flashcard.EaseFactor,
		Interval:     flashcard.Interval,
	}
}

func domainToModelSlice(flashcards []*domain.Flashcard) []*Flashcard {
	var models []*Flashcard
	for _, flashcard := range flashcards {
		models = append(models, domainToModel(flashcard))
	}
	return models
}

func modelToDomainSlice(flashcards []*Flashcard) []*domain.Flashcard {
	var domains []*domain.Flashcard
	for _, flashcard := range flashcards {
		domains = append(domains, flashcard.ToDomain())
	}
	return domains
}
