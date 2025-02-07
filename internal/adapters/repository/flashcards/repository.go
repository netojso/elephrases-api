package repository

import (
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
)

type flashcardRepository struct {
	DB *gorm.DB
}

func NewFlashcardRepository(db *gorm.DB) portrepository.FlashcardRepository {
	return &flashcardRepository{DB: db}
}

func (r *flashcardRepository) FindAll(options *portrepository.Options) ([]*domain.Flashcard, error) {
	flashcards := []*Flashcard{}

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

	return modelToDomainSlice(flashcards), nil
}

func (r *flashcardRepository) FindByID(id string) (*domain.Flashcard, error) {
	flashcard := Flashcard{}
	err := r.DB.Where("id = ?", id).First(&flashcard).Error
	if err != nil {
		return nil, err
	}
	return flashcard.ToDomain(), nil
}

func (r *flashcardRepository) Save(flashcard *domain.Flashcard) error {

	model := domainToModel(flashcard)

	err := r.DB.Create(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *flashcardRepository) Update(flashcard *domain.Flashcard) error {
	model := domainToModel(flashcard)

	err := r.DB.Save(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *flashcardRepository) Delete(id string) error {
	err := r.DB.Where("id = ?", id).Delete(&Flashcard{}).Error
	if err != nil {
		return err
	}
	return nil
}
