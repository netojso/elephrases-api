package usecase

import "github.com/netojso/elephrases-api/domain"

type FlashcardUsecase struct {
	flashcardRepo domain.FlashcardRepository
}

func NewFlashcardUsecase(flashcardRepo domain.FlashcardRepository) domain.FlashcardUsecase {
	return &FlashcardUsecase{
		flashcardRepo: flashcardRepo,
	}
}

func (fu *FlashcardUsecase) GetAll() ([]domain.Flashcard, error) {
	return fu.flashcardRepo.FindAll()
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

func (fu *FlashcardUsecase) Delete(id string) error {
	return fu.flashcardRepo.Delete(id)
}
