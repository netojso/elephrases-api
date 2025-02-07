package repository

import (
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
)

type DeckRepository struct {
	DB *gorm.DB
}

func NewDeckRepository(db *gorm.DB) portrepository.DeckRepository {
	return &DeckRepository{DB: db}
}

func (r *DeckRepository) FindAll() ([]*domain.Deck, error) {
	decks := []*Deck{}

	err := r.DB.Preload("Flashcards").Find(&decks).Error
	if err != nil {
		return nil, err
	}

	return modelToDomainSlice(decks), nil
}

func (r *DeckRepository) FindByID(id string) (*domain.Deck, error) {
	deck := Deck{}
	err := r.DB.Where("id = ?", id).Preload("Flashcards").First(&deck).Error
	if err != nil {
		return nil, err
	}
	return deck.ToDomain(), nil
}

func (r *DeckRepository) Save(deck *domain.Deck) error {
	model := domainToModel(deck)
	err := r.DB.Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepository) Update(deck *domain.Deck) error {
	model := domainToModel(deck)
	err := r.DB.Save(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepository) Delete(id string) error {
	err := r.DB.Delete(&Deck{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
