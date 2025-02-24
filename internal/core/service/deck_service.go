package service

import (
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
	decks, err := ds.repo.FindAll()

	if err != nil {
		return nil, err
	}

	for _, deck := range decks {

		deck.Stats.LearningCards = 0
		deck.Stats.ReviewingCards = 0
		deck.Stats.NewCards = 0

		for _, flashcard := range deck.Flashcards {
			switch flashcard.State {
			case "new":
				deck.Stats.NewCards++
			case "learning":
				deck.Stats.LearningCards++
			case "review":
				deck.Stats.ReviewingCards++
			}
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
