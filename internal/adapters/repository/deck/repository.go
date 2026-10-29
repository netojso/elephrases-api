package repository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
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

func (r *DeckRepository) FindAll(ctx context.Context) ([]*domain.Deck, error) {
	decks := []*Deck{}

	err := r.DB.WithContext(ctx).Preload("Flashcards").Find(&decks).Error
	if err != nil {
		return nil, err
	}

	return modelToDomainSlice(decks), nil
}

func (r *DeckRepository) FindByID(ctx context.Context, id string) (*domain.Deck, error) {
	deck := Deck{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).Preload("Flashcards").First(&deck).Error
	if err != nil {
		return nil, gormx.MapError(err)
	}
	return deck.ToDomain(), nil
}

func (r *DeckRepository) Save(ctx context.Context, deck *domain.Deck) error {
	model := domainToModel(deck)
	err := r.DB.WithContext(ctx).Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepository) Update(ctx context.Context, deck *domain.Deck) error {
	model := domainToModel(deck)
	err := r.DB.WithContext(ctx).Save(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepository) Delete(ctx context.Context, id string) error {
	return gormx.DeleteResult(r.DB.WithContext(ctx).Where("id = ?", id).Delete(&Deck{}))
}
