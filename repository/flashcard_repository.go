package repository

import (
	"github.com/netojso/elephrases-api/domain"
	"gorm.io/gorm"
)

type FlashcardRepository struct {
	DB *gorm.DB
}

func NewFlashcardRepository(db *gorm.DB) domain.FlashcardRepository {
	return &FlashcardRepository{DB: db}
}

func (r *FlashcardRepository) FindAll(options *domain.Options) ([]domain.Flashcard, error) {
	flashcards := []domain.Flashcard{}

	db := r.DB

	if options != nil {
		for key, value := range options.Where {
			db = db.Where(key, value)
		}
	}

	err := db.Find(&flashcards).Error

	if err != nil {
		return nil, err
	}

	return flashcards, nil
}

func (r *FlashcardRepository) FindByID(id string) (domain.Flashcard, error) {
	flashcard := domain.Flashcard{}
	err := r.DB.Where("id = ?", id).First(&flashcard).Error
	if err != nil {
		return domain.Flashcard{}, err
	}
	return flashcard, nil
}

func (r *FlashcardRepository) Save(flashcard domain.Flashcard) error {
	err := r.DB.Create(&flashcard).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FlashcardRepository) Update(flashcard domain.Flashcard) error {
	err := r.DB.Save(&flashcard).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FlashcardRepository) Delete(id string) error {
	err := r.DB.Where("id = ?", id).Delete(&domain.Flashcard{}).Error
	if err != nil {
		return err
	}
	return nil
}
