package repository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type flashcardRepository struct {
	DB *gorm.DB
}

func NewFlashcardRepository(db *gorm.DB) portrepository.FlashcardRepository {
	return &flashcardRepository{DB: db}
}

func (r *flashcardRepository) FindAll(ctx context.Context, options *portrepository.Options) ([]*domain.Flashcard, error) {
	flashcards := []*Flashcard{}

	db := gormx.DB(ctx, r.DB)

	if options != nil {
		for key, value := range options.Where {
			if value == nil {
				db = db.Where(key)
				continue
			}
			db = db.Where(key, value)
		}

		if options.Limit > 0 {
			db = db.Limit(options.Limit)
		}
	}

	err := db.Find(&flashcards).Error

	if err != nil {
		return nil, err
	}

	return modelToDomainSlice(flashcards), nil
}

func (r *flashcardRepository) FindByID(ctx context.Context, id string) (*domain.Flashcard, error) {
	flashcard := Flashcard{}
	err := gormx.DB(ctx, r.DB).Where("id = ?", id).First(&flashcard).Error
	if err != nil {
		return nil, gormx.MapError(err)
	}
	return flashcard.ToDomain(), nil
}

func (r *flashcardRepository) FindByIDForUpdate(ctx context.Context, id string) (*domain.Flashcard, error) {
	flashcard := Flashcard{}
	err := gormx.DB(ctx, r.DB).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&flashcard).Error
	if err != nil {
		return nil, gormx.MapError(err)
	}
	return flashcard.ToDomain(), nil
}

func (r *flashcardRepository) FindByDeckID(ctx context.Context, deckID string) ([]*domain.Flashcard, error) {
	flashcards := []*Flashcard{}
	err := gormx.DB(ctx, r.DB).Where("deck_id = ?", deckID).Find(&flashcards).Error
	if err != nil {
		return nil, err
	}

	return modelToDomainSlice(flashcards), nil
}

func (r *flashcardRepository) Save(ctx context.Context, flashcard *domain.Flashcard) error {

	model := domainToModel(flashcard)

	err := gormx.DB(ctx, r.DB).Create(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *flashcardRepository) Update(ctx context.Context, flashcard *domain.Flashcard) error {

	model := domainToModel(flashcard)

	err := gormx.DB(ctx, r.DB).Save(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *flashcardRepository) Delete(ctx context.Context, id string) error {
	return gormx.DeleteResult(gormx.DB(ctx, r.DB).Where("id = ?", id).Delete(&Flashcard{}))
}
