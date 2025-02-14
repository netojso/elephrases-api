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

// GetAll godoc
// @Summary Get all decks
// @Description Get all decks
// @Tags Decks
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} domain.Deck
// @Failure 404 {object} object{error=string}
// @Router /decks [get]
func (dh *DeckHandler) GetAll(ctx *gin.Context) {
	decks, err := dh.deckService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, decks)
}

// GetByID godoc
// @Summary Get a deck by ID
// @Description Get a deck by ID
// @Tags Decks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 200 {object} domain.Deck
// @Failure 404 {object} object{error=string}
// @Router /decks/{id} [get]
func (dh *DeckHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	deck, err := dh.deckService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, deck)
}

// Create godoc
// @Summary Create a new deck
// @Description Create a new deck
// @Tags Decks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param deck body dto.CreateDeckDTO true "Create Deck DTO"
// @Success 201
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /decks [post]
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

	ctx.JSON(http.StatusCreated, nil)
}

// Update godoc
// @Summary Update a deck
// @Description Update a deck
// @Tags Decks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Param deck body domain.Deck true "Update Deck"
// @Success 200 {object} domain.Deck
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /decks/{id} [put]
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

// Delete godoc
// @Summary Delete a deck
// @Description Delete a deck
// @Tags Decks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 204
// @Failure 500 {object} object{error=string}
// @Router /decks/{id} [delete]
func (dh *DeckHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := dh.deckService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
