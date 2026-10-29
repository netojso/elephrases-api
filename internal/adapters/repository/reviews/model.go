package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/netojso/elephrases-api/internal/core/domain"
)

type Review struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	FlashcardID   uuid.UUID `gorm:"type:uuid"`
	Response      string    `gorm:"type:varchar(20)"`
	PreviousState string    `gorm:"type:varchar(20)"`
	NewState      string    `gorm:"type:varchar(20)"`
	IntervalNs    int64     `gorm:"column:interval_ns;type:bigint"`
	EaseFactor    float64   `gorm:"type:float"`
	ReviewedAt    time.Time `gorm:"type:timestamp"`
}

func (r *Review) TableName() string {
	return "reviews"
}

func domainToModel(review *domain.Review) *Review {
	return &Review{
		ID:            review.ID.Value(),
		FlashcardID:   review.FlashcardID.Value(),
		Response:      review.Response,
		PreviousState: string(review.PreviousState),
		NewState:      string(review.NewState),
		IntervalNs:    int64(review.Interval),
		EaseFactor:    review.EaseFactor,
		ReviewedAt:    review.ReviewedAt,
	}
}
