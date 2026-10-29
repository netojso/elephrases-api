package service

import (
	"context"
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg/nullable"
	"github.com/spf13/viper"
)

type flashcardService struct {
	repo    portrepository.FlashcardRepository
	reviews portrepository.ReviewRepository
	tx      portrepository.TxManager
}

func NewFlashcardService(
	repo portrepository.FlashcardRepository,
	reviews portrepository.ReviewRepository,
	tx portrepository.TxManager,
) portservice.FlashcardService {
	return &flashcardService{repo: repo, reviews: reviews, tx: tx}
}

func (fs *flashcardService) Study(ctx context.Context, deckID string) (*domain.Study, error) {
	dailyLimit := 20

	yesterday := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0,
		0,
		0,
		0,
		time.Now().Location(),
	)

	review_today_flashcards, err := fs.repo.FindAll(
		ctx,
		&portrepository.Options{
			Where: map[string]interface{}{
				"deck_id = ?":        deckID,
				"last_review_at > ?": yesterday,
			},
		})

	if err != nil {
		return nil, err
	}

	new_flashcards, err := fs.repo.FindAll(
		ctx,
		&portrepository.Options{
			Where: map[string]interface{}{
				"deck_id = ?":            deckID,
				"last_review_at IS NULL": nil,
			},
			Limit: dailyLimit - len(review_today_flashcards),
		},
	)

	if err != nil {
		return nil, err
	}

	due_flashcards, err := fs.repo.FindAll(
		ctx,
		&portrepository.Options{
			Where: map[string]interface{}{
				"deck_id = ?":         deckID,
				"next_review_at <= ?": time.Now(),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if new_flashcards == nil {
		new_flashcards = []*domain.Flashcard{}
	}

	if due_flashcards == nil {
		due_flashcards = []*domain.Flashcard{}
	}

	return &domain.Study{
		NewCards: new_flashcards,
		DueCards: due_flashcards,
	}, nil
}

func (fs *flashcardService) GetDueFlashcards(ctx context.Context) ([]*domain.Flashcard, error) {
	return fs.repo.FindAll(
		ctx,
		&portrepository.Options{
			Where: map[string]interface{}{
				"next_review_at <= ?": time.Now(),
			},
		},
	)
}

func (fs *flashcardService) GetAll(ctx context.Context) ([]*domain.Flashcard, error) {
	return fs.repo.FindAll(ctx, nil)
}

func (fs *flashcardService) GetByID(ctx context.Context, id string) (*domain.Flashcard, error) {
	return fs.repo.FindByID(ctx, id)
}

func (fs *flashcardService) GetByDeckID(ctx context.Context, deckID string) ([]*domain.Flashcard, error) {
	return fs.repo.FindByDeckID(ctx, deckID)
}

func (fs *flashcardService) Create(ctx context.Context, flashcard *domain.Flashcard) error {
	// add s3 domain to media flashcard

	s3Url := viper.GetString("AWS_S3_URL")

	mediaUrl := nullable.NewNullableString("")

	if flashcard.MediaUrl.Valid {
		mediaUrl = nullable.NewNullableString(s3Url + flashcard.MediaUrl.String)
	}

	flashcard.MediaUrl = mediaUrl

	return fs.repo.Save(ctx, flashcard)
}

func (fs *flashcardService) Update(ctx context.Context, flashcard *domain.Flashcard) error {
	return fs.repo.Update(ctx, flashcard)
}

// Review aplica a resposta ao cartão e registra a revisão na mesma transação.
// O cartão é bloqueado (FOR UPDATE) para que revisões concorrentes não se sobrescrevam.
func (fs *flashcardService) Review(ctx context.Context, id string, response string) error {
	return fs.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		flashcard, err := fs.repo.FindByIDForUpdate(ctx, id)

		if err != nil {
			return err
		}

		previousState := flashcard.State

		flashcard.ReviewFlashcard(response, nil)

		if err := fs.repo.Update(ctx, flashcard); err != nil {
			return err
		}

		return fs.reviews.Save(ctx, domain.NewReview(flashcard, response, previousState, flashcard.LastReviewAt.Time))
	})
}

func (fs *flashcardService) Delete(ctx context.Context, id string) error {
	return fs.repo.Delete(ctx, id)
}
