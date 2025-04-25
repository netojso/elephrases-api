package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/dto"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg"
)

type FlashcardHandler struct {
	service portservice.FlashcardService
	storage portrepository.StoragePort
}

func NewFlashcardHandler(
	service portservice.FlashcardService,
	storage portrepository.StoragePort,
) *FlashcardHandler {
	return &FlashcardHandler{
		service: service,
		storage: storage,
	}
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

func (fh FlashcardHandler) Study(ctx *gin.Context) {
	deckID := ctx.Param("deckID")

	data, err := fh.service.Study(deckID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
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

// GetByDeckID godoc
// @Summary Get flashcards by Deck ID
// @Description Get all flashcards by Deck ID
// @Tags Flashcards
// @Produce json
// @Security BearerAuth
// @Param deckID path string true "Deck ID"
// @Success 200 {array} domain.Flashcard
// @Failure 404 {object} object{error=string}
// @Router /flashcards/deck/{deckID} [get]
func (fh FlashcardHandler) GetByDeckID(ctx *gin.Context) {
	deckID := ctx.Param("deckID")

	flashcards, err := fh.service.GetByDeckID(deckID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flashcards)
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

	flashcard := domain.NewFlashcard(DeckID, body.Front, body.Back, body.Media)

	if err := fh.service.Create(flashcard); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, flashcard)
}

func (fh FlashcardHandler) CreateMany(ctx *gin.Context) {
	flashcardsJSON := ctx.PostForm("flashcards")

	var flashcards []dto.CreateFlashcardDTO
	if err := json.Unmarshal([]byte(flashcardsJSON), &flashcards); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	form, _ := ctx.MultipartForm()
	files := form.File["media"]

	for _, flashcard := range flashcards {
		mediaName := flashcard.Media // mediaName is a pointer to a string
		mediaData := []byte{}

		if mediaName != "" {
			for _, file := range files {
				if file.Filename == mediaName {
					fileContent, err := file.Open()
					if err != nil {
						ctx.JSON(500, gin.H{"error": "Error opening file"})
						return
					}

					defer fileContent.Close()

					fileBytes, err := io.ReadAll(fileContent)
					if err != nil {
						ctx.JSON(500, gin.H{"error": "Error reading file"})
						return
					}

					mediaData = fileBytes
					break
				}
			}

			if len(mediaData) == 0 {
				ctx.JSON(400, gin.H{"error": "Media file not found"})
				return
			}

			err := fh.storage.Upload(portrepository.File{
				Name: mediaName,
				Data: mediaData,
			})

			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		DeckID, err := pkg.ParseUUID(flashcard.DeckID)

		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		newFlashcard := domain.NewFlashcard(DeckID, flashcard.Front, flashcard.Back, mediaName)

		if err := fh.service.Create(newFlashcard); err != nil {
			ctx.JSON(500, gin.H{"error": "Error creating flashcard"})
			return
		}
	}

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
