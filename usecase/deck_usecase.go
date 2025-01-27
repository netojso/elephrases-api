package usecase

import "github.com/netojso/elephrases-api/domain"

type DeckUsecase struct {
	deckRepo domain.DeckRepository
}

func NewDeckUsecase(deckRepo domain.DeckRepository) domain.DeckUsecase {
	return &DeckUsecase{
		deckRepo: deckRepo,
	}
}

func (fu *DeckUsecase) GetAll() ([]domain.Deck, error) {
	return fu.deckRepo.FindAll()
}

func (fu *DeckUsecase) GetByID(id string) (domain.Deck, error) {
	return fu.deckRepo.FindByID(id)
}

func (fu *DeckUsecase) Create(deck domain.Deck) error {
	return fu.deckRepo.Save(deck)
}

func (fu *DeckUsecase) Update(deck domain.Deck) error {
	return fu.deckRepo.Update(deck)
}

func (fu *DeckUsecase) Delete(id string) error {
	return fu.deckRepo.Delete(id)
}
