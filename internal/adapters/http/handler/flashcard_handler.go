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

// GetDueFlashcards godoc
// @Summary Get due flashcards
// @Description Get all due flashcards
// @Tags Flashcards
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Flashcard
// @Failure 500 {object} object{error=string}
// @Router /flashcards/due [get]
func (fh FlashcardHandler) GetDueFlashcards(ctx *gin.Context) {
	flashcards, err := fh.service.GetDueFlashcards()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
}

// Review godoc
// @Summary Review a flashcard
// @Description Review a flashcard by providing the flashcard ID and response
// @Tags Flashcards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body dto.ReviewFlashcardDTO true "Review Flashcard"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /flashcards/review [post]
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

// GetAll godoc
// @Summary Get all flashcards
// @Description Get all flashcards
// @Tags Flashcards
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Flashcard
// @Failure 500 {object} object{error=string}
// @Router /flashcards [get]
func (fh FlashcardHandler) GetAll(ctx *gin.Context) {
	flashcards, err := fh.service.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
}

// GetByID godoc
// @Summary Get a flashcard by ID
// @Description Get a flashcard by its ID
// @Tags Flashcards
// @Produce json
// @Security BearerAuth
// @Param id path string true "Flashcard ID"
// @Success 200 {object} domain.Flashcard
// @Failure 404 {object} object{error=string}
// @Router /flashcards/{id} [get]
func (fh FlashcardHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	flashcard, err := fh.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcard)
}

// Create godoc
// @Summary Create a new flashcard
// @Description Create a new flashcard with the provided details
// @Tags Flashcards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param flashcard body dto.CreateFlashcardDTO true "Create Flashcard"
// @Success 201 {object} domain.Flashcard
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /flashcards [post]
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

// Update godoc
// @Summary Update a flashcard
// @Description Update a flashcard by its ID
// @Tags Flashcards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Flashcard ID"
// @Param flashcard body domain.Flashcard true "Update Flashcard"
// @Success 200 {object} domain.Flashcard
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /flashcards/{id} [put]
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

// Delete godoc
// @Summary Delete a flashcard
// @Description Delete a flashcard by its ID
// @Tags Flashcards
// @Produce json
// @Security BearerAuth
// @Param id path string true "Flashcard ID"
// @Success 204
// @Failure 500 {object} object{error=string}
// @Router /flashcards/{id} [delete]
func (fh FlashcardHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := fh.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
