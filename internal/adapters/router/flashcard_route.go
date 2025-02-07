package router

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/handler"
	repository "github.com/netojso/elephrases-api/internal/adapters/repository/flashcards"
	"github.com/netojso/elephrases-api/internal/core/service"
	"gorm.io/gorm"
)

func NewFlashcardRouter(db *gorm.DB, group *gin.RouterGroup) {
	repo := repository.NewFlashcardRepository(db)
	service := service.NewFlashcardService(repo)
	handler := handler.NewFlashcardHandler(service)

	group.GET("/flashcards/due", handler.GetDueFlashcards)

	group.POST("/flashcards/review", handler.Review)

	group.GET("/flashcards", handler.GetAll)

	group.GET("/flashcards/:id", handler.GetByID)

	group.POST("/flashcards", handler.Create)

	group.PUT("/flashcards/:id", handler.Update)

	group.DELETE("/flashcards/:id", handler.Delete)

}
