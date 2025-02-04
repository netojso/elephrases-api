package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/domain"
)

type FlashcardController struct {
	FlashcardUsecase domain.FlashcardUsecase
}

func (fc *FlashcardController) GetDueFlashcards(ctx *gin.Context) {
	flashcards, err := fc.FlashcardUsecase.GetDueFlashcards()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
}

func (fc *FlashcardController) Review(ctx *gin.Context) {

	var body struct {
		FlashCardID string `json:"flashcard_id" binding:"required"`
		Response    string `json:"response" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.FlashcardUsecase.Review(body.FlashCardID, body.Response); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (fc *FlashcardController) GetAll(ctx *gin.Context) {

	flashcards, err := fc.FlashcardUsecase.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
}

func (fc *FlashcardController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	flashcard, err := fc.FlashcardUsecase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcard)

}

func (fc *FlashcardController) Create(ctx *gin.Context) {
	var body domain.CreateFlashcardRequest

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flashcard, err := domain.NewFlashcard(body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.FlashcardUsecase.Create(*flashcard); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, flashcard)
}

func (fc *FlashcardController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	flashcard, err := fc.FlashcardUsecase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBind(&flashcard); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.FlashcardUsecase.Update(flashcard); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcard)
}

func (fc *FlashcardController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := fc.FlashcardUsecase.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
