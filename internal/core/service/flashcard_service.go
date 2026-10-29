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
	now     func() time.Time
}

// FlashcardServiceOption personaliza o flashcardService.
type FlashcardServiceOption func(*flashcardService)

// WithClock substitui time.Now; usado nos testes para tornar o tempo determinístico.
func WithClock(now func() time.Time) FlashcardServiceOption {
	return func(fs *flashcardService) {
		fs.now = now
	}
}

func NewFlashcardService(
	repo portrepository.FlashcardRepository,
	reviews portrepository.ReviewRepository,
	tx portrepository.TxManager,
	opts ...FlashcardServiceOption,
) portservice.FlashcardService {
	fs := &flashcardService{repo: repo, reviews: reviews, tx: tx, now: time.Now}
	for _, opt := range opts {
		opt(fs)
	}
	return fs
}

// dailyNewCardLimit é o máximo de cartões estudados por deck por dia.
const dailyNewCardLimit = 20

func (fs *flashcardService) Study(ctx context.Context, deckID string) (*domain.Study, error) {
	return fs.repo.FindStudyQueue(ctx, deckID, dailyNewCardLimit, fs.now())
}

func (fs *flashcardService) GetDueFlashcards(ctx context.Context) ([]*domain.Flashcard, error) {
	return fs.repo.FindAll(
		ctx,
		&portrepository.Options{
			Where: map[string]interface{}{
				"next_review_at <= ?": fs.now(),
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
		now := fs.now()

		flashcard.ReviewFlashcard(response, nil, now)

		if err := fs.repo.Update(ctx, flashcard); err != nil {
			return err
		}

		return fs.reviews.Save(ctx, domain.NewReview(flashcard, response, previousState, now))
	})
}

func (fs *flashcardService) Delete(ctx context.Context, id string) error {
	return fs.repo.Delete(ctx, id)
}
