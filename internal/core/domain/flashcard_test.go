package domain

import (
	"testing"
	"time"

	"github.com/netojso/elephrases-api/pkg/nullable"
	"github.com/stretchr/testify/assert"
)

var reviewNow = time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)

// Configuração padrão de ReviewFlashcard: passos de 1m e 10m, graduação em 1d, easy em 4d.
const (
	step1      = time.Minute
	step2      = 10 * time.Minute
	graduating = 24 * time.Hour
	easy       = 4 * 24 * time.Hour
)

func TestReviewFlashcard(t *testing.T) {
	previousNext := reviewNow.Add(-time.Hour) // cartão em review já vencido

	tests := []struct {
		name         string
		state        CardState
		interval     time.Duration
		response     string
		wantState    CardState
		wantInterval time.Duration
		// wantNext é o next_review_at esperado; zero significa "não reagendado".
		wantNext time.Time
	}{
		// new: entra em learning no primeiro passo antes de aplicar a resposta.
		{"new/again", StateNew, 0, "again", StateLearning, step1, reviewNow.Add(step1)},
		{"new/hard", StateNew, 0, "hard", StateLearning, (step1 + step2) / 2, reviewNow.Add((step1 + step2) / 2)},
		{"new/good", StateNew, 0, "good", StateLearning, step2, reviewNow.Add(step2)},
		{"new/easy", StateNew, 0, "easy", StateReview, easy, reviewNow.Add(easy)},

		// learning no primeiro passo.
		{"learning step1/again", StateLearning, step1, "again", StateLearning, step1, reviewNow.Add(step1)},
		{"learning step1/hard", StateLearning, step1, "hard", StateLearning, (step1 + step2) / 2, reviewNow.Add((step1 + step2) / 2)},
		{"learning step1/good", StateLearning, step1, "good", StateLearning, step2, reviewNow.Add(step2)},
		{"learning step1/easy", StateLearning, step1, "easy", StateReview, easy, reviewNow.Add(easy)},

		// learning no último passo: good gradua.
		{"learning step2/again", StateLearning, step2, "again", StateLearning, step1, reviewNow.Add(step1)},
		{"learning step2/hard", StateLearning, step2, "hard", StateLearning, step2, reviewNow.Add(step2)},
		{"learning step2/good", StateLearning, step2, "good", StateReview, graduating, reviewNow.Add(graduating)},
		{"learning step2/easy", StateLearning, step2, "easy", StateReview, easy, reviewNow.Add(easy)},

		// review.
		{"review/again", StateReview, graduating, "again", StateReview, step1, reviewNow.Add(step1)},
		// BUG CONHECIDO: hard e good em review não reagendam o cartão; next_review_at
		// continua no passado e o cartão fica vencido para sempre. Corrigir exige
		// implementar o SM-2 (ease_factor e multiplicação do intervalo).
		{"review/hard (bug conhecido)", StateReview, graduating, "hard", StateReview, graduating, previousNext},
		{"review/good (bug conhecido)", StateReview, graduating, "good", StateReview, graduating, previousNext},
		{"review/easy", StateReview, graduating, "easy", StateReview, easy, reviewNow.Add(easy)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flashcard{State: tt.state, Interval: tt.interval, EaseFactor: DefaultEaseFactor}
			if tt.state == StateReview {
				f.NextReviewAt = nullable.NewNullableTime(previousNext)
			}

			f.ReviewFlashcard(tt.response, nil, reviewNow)

			assert.Equal(t, tt.wantState, f.State, "state")
			assert.Equal(t, tt.wantInterval, f.Interval, "interval")
			assert.True(t, f.NextReviewAt.Valid, "next_review_at deve estar preenchido")
			assert.True(t, tt.wantNext.Equal(f.NextReviewAt.Time), "next_review_at: want %s, got %s", tt.wantNext, f.NextReviewAt.Time)
			assert.True(t, reviewNow.Equal(f.LastReviewAt.Time), "last_review_at deve ser o relógio injetado")
			// O algoritmo atual não ajusta o ease_factor em nenhuma resposta.
			assert.Equal(t, DefaultEaseFactor, f.EaseFactor, "ease_factor")
		})
	}
}

func TestReviewFlashcard_UnknownResponseReschedulesWithCurrentInterval(t *testing.T) {
	f := &Flashcard{State: StateLearning, Interval: step2, EaseFactor: DefaultEaseFactor}

	f.ReviewFlashcard("xyz", nil, reviewNow)

	assert.Equal(t, StateLearning, f.State)
	assert.True(t, reviewNow.Add(step2).Equal(f.NextReviewAt.Time))
}

func TestReviewFlashcard_CustomSettings(t *testing.T) {
	settings := &Settings{
		LearningSteps:      []time.Duration{5 * time.Minute},
		GraduatingInterval: 48 * time.Hour,
		EasyInterval:       7 * 24 * time.Hour,
	}

	t.Run("passo único: good gradua", func(t *testing.T) {
		f := &Flashcard{State: StateNew, EaseFactor: DefaultEaseFactor}
		f.ReviewFlashcard("good", settings, reviewNow)
		assert.Equal(t, StateReview, f.State)
		assert.Equal(t, 48*time.Hour, f.Interval)
	})

	t.Run("passo único: hard multiplica por 1.5", func(t *testing.T) {
		f := &Flashcard{State: StateNew, EaseFactor: DefaultEaseFactor}
		f.ReviewFlashcard("hard", settings, reviewNow)
		assert.Equal(t, StateLearning, f.State)
		assert.Equal(t, 7*time.Minute+30*time.Second, f.Interval)
	})
}

func TestNewFlashcard(t *testing.T) {
	f := NewFlashcard(NewDeck("d", "", "c", "v").ID, "front", "back", "")

	assert.Equal(t, StateNew, f.State)
	assert.Equal(t, DefaultEaseFactor, f.EaseFactor)
	assert.False(t, f.MediaUrl.Valid)
	assert.False(t, f.LastReviewAt.Valid)
}
