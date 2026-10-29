package repository

import (
	"context"
	"time"

	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
	"github.com/netojso/elephrases-api/internal/core/domain"
)

const (
	queueNew = 0
	queueDue = 1
)

// studyQueueSQL monta a fila de estudo em uma única ida ao banco:
//   - reviewed_today: cartões do deck revisados desde o início do dia;
//   - new_cards: nunca revisados, até o limite diário restante (nunca negativo);
//   - due_cards: vencidos (next_review_at <= now).
//
// Usa os índices (deck_id, last_review_at) e (deck_id, next_review_at) da 000007.
const studyQueueSQL = `
WITH reviewed_today AS (
    SELECT COUNT(*) AS total
    FROM flashcards
    WHERE deck_id = ? AND last_review_at >= ?
),
new_cards AS (
    SELECT f.*, 0 AS queue
    FROM flashcards f
    WHERE f.deck_id = ? AND f.last_review_at IS NULL
    ORDER BY f.created_at, f.id
    LIMIT GREATEST(?::int - (SELECT total FROM reviewed_today), 0)
),
due_cards AS (
    SELECT f.*, 1 AS queue
    FROM flashcards f
    WHERE f.deck_id = ? AND f.next_review_at <= ?
)
SELECT * FROM new_cards
UNION ALL
SELECT * FROM due_cards
ORDER BY queue, next_review_at, created_at, id`

type studyRow struct {
	Flashcard `gorm:"embedded"`
	Queue     int
}

func (r *flashcardRepository) FindStudyQueue(ctx context.Context, deckID string, dailyLimit int, now time.Time) (*domain.Study, error) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var rows []studyRow
	err := gormx.DB(ctx, r.DB).
		Raw(studyQueueSQL, deckID, startOfDay, deckID, dailyLimit, deckID, now).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	study := &domain.Study{
		NewCards: []*domain.Flashcard{},
		DueCards: []*domain.Flashcard{},
	}

	for i := range rows {
		flashcard := rows[i].Flashcard.ToDomain()
		switch rows[i].Queue {
		case queueNew:
			study.NewCards = append(study.NewCards, flashcard)
		case queueDue:
			study.DueCards = append(study.DueCards, flashcard)
		}
	}

	return study, nil
}
