package usecase

import (
	"time"

	"github.com/netojso/elephrases-api/domain"
)

type FlashcardUsecase struct {
	flashcardRepo domain.FlashcardRepository
}

func NewFlashcardUsecase(flashcardRepo domain.FlashcardRepository) domain.FlashcardUsecase {
	return &FlashcardUsecase{
		flashcardRepo: flashcardRepo,
	}
}

func (fu *FlashcardUsecase) GetDueFlashcards() ([]domain.Flashcard, error) {
	return fu.flashcardRepo.FindAll(
		&domain.Options{
			Where: map[string]interface{}{
				"next_review_at <= ?": time.Now(),
			},
		},
	)
}

func (fu *FlashcardUsecase) GetAll() ([]domain.Flashcard, error) {
	return fu.flashcardRepo.FindAll(nil)
}

func (fu *FlashcardUsecase) GetByID(id string) (domain.Flashcard, error) {
	return fu.flashcardRepo.FindByID(id)
}

func (fu *FlashcardUsecase) Create(flashcard domain.Flashcard) error {
	return fu.flashcardRepo.Save(flashcard)
}

func (fu *FlashcardUsecase) Update(flashcard domain.Flashcard) error {
	return fu.flashcardRepo.Update(flashcard)
}

func (fu *FlashcardUsecase) Review(id string, response string) error {
	flashcard, err := fu.flashcardRepo.FindByID(id)

	if err != nil {
		return err
	}

	flashcard.ReviewFlashcard(response, nil)

	return fu.flashcardRepo.Update(flashcard)
}

func (fu *FlashcardUsecase) Delete(id string) error {
	return fu.flashcardRepo.Delete(id)
}
