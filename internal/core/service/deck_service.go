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
	return ds.repo.FindAll()
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
