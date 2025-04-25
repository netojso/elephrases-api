package service

import (
	"time"

	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
)

type DeckService struct {
	repo portrepository.DeckRepository
}

func NewDeckUsecase(repo portrepository.DeckRepository) portservice.DeckService {
	return &DeckService{
		repo: repo,
	}
}

func (ds *DeckService) GetAll() ([]*domain.Deck, error) {
	dailyLimit := 20
	decks, err := ds.repo.FindAll()

	if err != nil {
		return nil, err
	}

	for _, deck := range decks {

		deck.Stats.LearningCards = 0
		deck.Stats.ReviewCards = 0
		deck.Stats.NewCards = 0
		deck.Stats.TotalCards = len(deck.Flashcards)

		totalReviewTotal := 0

		for _, flashcard := range deck.Flashcards {
			switch flashcard.State {
			case "new":
				deck.Stats.NewCards++
			case "learning":
				deck.Stats.LearningCards++
			case "review":
				deck.Stats.ReviewCards++
			}

			if flashcard.LastReviewAt.Valid {
				yesterday := time.Now().Add(-24 * time.Hour)
				if flashcard.LastReviewAt.Time.After(yesterday) {
					totalReviewTotal++
				}
			}
		}

		if deck.Stats.NewCards > dailyLimit {
			deck.Stats.NewCards = dailyLimit - totalReviewTotal
		}
	}

	return decks, nil
}

func (ds *DeckService) GetByID(id string) (*domain.Deck, error) {
	return ds.repo.FindByID(id)
}

func (ds *DeckService) Create(deck *domain.Deck) error {
	return ds.repo.Save(deck)
}

func (ds *DeckService) Update(deck *domain.Deck) error {
	return ds.repo.Update(deck)
}

func (ds *DeckService) Delete(id string) error {
	return ds.repo.Delete(id)
}
