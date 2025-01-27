package repository

import (
	"github.com/netojso/elephrases-api/domain"
	"gorm.io/gorm"
)

type DeckRepository struct {
	DB *gorm.DB
}

func NewDeckRepository(db *gorm.DB) domain.DeckRepository {
	return &DeckRepository{DB: db}
}

func (r *DeckRepository) FindAll() ([]domain.Deck, error) {
	decks := []domain.Deck{}
	err := r.DB.Preload("Flashcards").Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func (r *DeckRepository) FindByID(id string) (domain.Deck, error) {
	deck := domain.Deck{}
	err := r.DB.Where("id = ?", id).First(&deck).Error
	if err != nil {
		return domain.Deck{}, err
	}
	return deck, nil
}

func (r *DeckRepository) Save(deck domain.Deck) error {
	err := r.DB.Create(&deck).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepository) Update(deck domain.Deck) error {
	err := r.DB.Save(&deck).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepository) Delete(id string) error {
	err := r.DB.Where("id = ?", id).Delete(&domain.Deck{}).Error
	if err != nil {
		return err
	}
	return nil
}
