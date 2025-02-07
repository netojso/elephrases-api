package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/dto"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
)

type DeckHandler struct {
	deckService portservice.DeckService
}

func NewDeckHandler(deckService portservice.DeckService) *DeckHandler {
	return &DeckHandler{
		deckService: deckService,
	}
}

func (dh *DeckHandler) GetAll(ctx *gin.Context) {
	decks, err := dh.deckService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, decks)
}

func (dh *DeckHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	deck, err := dh.deckService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, deck)
}

func (dh *DeckHandler) Create(ctx *gin.Context) {
	var body *dto.CreateDeckDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deck := domain.NewDeck(body.Name, body.Description, body.Category, body.Visibility)

	if err := dh.deckService.Create(deck); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, deck)
}

func (dh *DeckHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	deck, err := dh.deckService.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&deck); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dh.deckService.Update(deck); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deck)
}

func (dh *DeckHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := dh.deckService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
