package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/domain"
)

type DeckController struct {
	DeckUsecase domain.DeckUsecase
}

func (fc *DeckController) GetAll(ctx *gin.Context) {

	decks, err := fc.DeckUsecase.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, decks)
}

func (fc *DeckController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	deck, err := fc.DeckUsecase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	ctx.JSON(http.StatusOK, deck)

}

func (fc *DeckController) Create(ctx *gin.Context) {
	var body domain.CreateDeckRequest

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deck, err := domain.NewDeck(body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.DeckUsecase.Create(*deck); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, deck)
}

func (fc *DeckController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	deck, err := fc.DeckUsecase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBind(&deck); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.DeckUsecase.Update(deck); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deck)
}

func (fc *DeckController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := fc.DeckUsecase.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
