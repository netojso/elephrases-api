package router

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/config"
	"github.com/netojso/elephrases-api/internal/adapters/http/handler"
	repository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/adapters/storage"
	"github.com/netojso/elephrases-api/internal/core/service"
	"gorm.io/gorm"
)

func NewFlashcardRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	repo := repository.NewFlashcardRepository(db)
	service := service.NewFlashcardService(repo)
	storage, _ := storage.NewS3Adapter(env)
	handler := handler.NewFlashcardHandler(service, storage)

	group.GET("/flashcards/deck/:deckID/study", handler.Study)

	group.GET("/flashcards/due", handler.GetDueFlashcards)

	group.POST("/flashcards/review", handler.Review)

	group.GET("/flashcards", handler.GetAll)

	group.GET("/flashcards/deck/:deckID", handler.GetByDeckID)

	group.GET("/flashcards/:id", handler.GetByID)

	group.POST("/flashcards", handler.Create)

	group.POST("/flashcards/create-many", handler.CreateMany)

	group.PUT("/flashcards/:id", handler.Update)

	group.DELETE("/flashcards/:id", handler.Delete)

}
