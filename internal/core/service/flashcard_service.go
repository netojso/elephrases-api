package service

import (
	"fmt"
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg/nullable"
	"github.com/spf13/viper"
)

type flashcardService struct {
	repo portrepository.FlashcardRepository
}

func NewFlashcardService(repo portrepository.FlashcardRepository) portservice.FlashcardService {
	return &flashcardService{repo: repo}
}

func (fs *flashcardService) Study(deckID string) (*domain.Study, error) {
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

	fmt.Println("yesterday", yesterday)

	review_today_flashcards, err := fs.repo.FindAll(
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
		&portrepository.Options{
			Where: map[string]interface{}{
				"deck_id":        deckID,
				"last_review_at": nil,
			},
			Limit: dailyLimit - len(review_today_flashcards),
		},
	)

	if err != nil {
		return nil, err
	}

	due_flashcards, err := fs.repo.FindAll(
		&portrepository.Options{
			Where: map[string]interface{}{
				"deck_id":             deckID,
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

func (fs *flashcardService) GetDueFlashcards() ([]*domain.Flashcard, error) {
	return fs.repo.FindAll(
		&portrepository.Options{
			Where: map[string]interface{}{
				"next_review_at <= ?": time.Now(),
			},
		},
	)
}

func (fs *flashcardService) GetAll() ([]*domain.Flashcard, error) {
	return fs.repo.FindAll(nil)
}

func (fs *flashcardService) GetByID(id string) (*domain.Flashcard, error) {
	return fs.repo.FindByID(id)
}

func (fs *flashcardService) GetByDeckID(deckID string) ([]*domain.Flashcard, error) {
	return fs.repo.FindByDeckID(deckID)
}

func (fs *flashcardService) Create(flashcard *domain.Flashcard) error {
	// add s3 domain to media flashcard

	s3Url := viper.GetString("AWS_S3_URL")

	mediaUrl := nullable.NewNullableString("")

	if flashcard.MediaUrl.Valid {
		mediaUrl = nullable.NewNullableString(s3Url + flashcard.MediaUrl.String)
	}

	flashcard.MediaUrl = mediaUrl

	return fs.repo.Save(flashcard)
}

func (fs *flashcardService) Update(flashcard *domain.Flashcard) error {
	return fs.repo.Update(flashcard)
}

func (fs *flashcardService) Review(id string, response string) error {
	flashcard, err := fs.repo.FindByID(id)

	if err != nil {
		return err
	}

	flashcard.ReviewFlashcard(response, nil)

	return fs.repo.Update(flashcard)
}

func (fs *flashcardService) Delete(id string) error {
	return fs.repo.Delete(id)
}
