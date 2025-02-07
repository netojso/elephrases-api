package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/dto"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg"
)

type FlashcardHandler struct {
	service portservice.FlashcardService
}

func NewFlashcardHandler(service portservice.FlashcardService) *FlashcardHandler {
	return &FlashcardHandler{service: service}
}

func (fh FlashcardHandler) GetDueFlashcards(ctx *gin.Context) {
	flashcards, err := fh.service.GetDueFlashcards()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
}

func (fh FlashcardHandler) Review(ctx *gin.Context) {
	var body dto.ReviewFlashcardDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fh.service.Review(body.FlashCardID, body.Response); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (fh FlashcardHandler) GetAll(ctx *gin.Context) {

	flashcards, err := fh.service.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
}

func (fh FlashcardHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	flashcard, err := fh.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcard)

}

func (fh FlashcardHandler) Create(ctx *gin.Context) {
	var body dto.CreateFlashcardDTO

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DeckID, err := pkg.ParseUUID(body.DeckID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flashcard := domain.NewFlashcard(DeckID, body.Front, body.Back)

	if err := fh.service.Create(flashcard); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, flashcard)
}

func (fh FlashcardHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	flashcard, err := fh.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBind(&flashcard); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fh.service.Update(flashcard); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcard)
}

func (fh FlashcardHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := fh.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
