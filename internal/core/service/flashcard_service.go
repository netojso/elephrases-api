package service

import (
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
)

type flashcardService struct {
	repo portrepository.FlashcardRepository
}

func NewFlashcardService(repo portrepository.FlashcardRepository) portservice.FlashcardService {
	return &flashcardService{repo: repo}
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

func (fs *flashcardService) Create(flashcard *domain.Flashcard) error {
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
