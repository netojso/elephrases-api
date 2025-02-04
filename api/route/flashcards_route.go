package route

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/api/controller"
	"github.com/netojso/elephrases-api/repository"
	"github.com/netojso/elephrases-api/usecase"
	"gorm.io/gorm"
)

func NewFlashcardRoute(db *gorm.DB, group *gin.RouterGroup) {
	fr := repository.NewFlashcardRepository(db)
	fc := &controller.FlashcardController{
		FlashcardUsecase: usecase.NewFlashcardUsecase(fr),
	}

	group.GET("/flashcards/due", fc.GetDueFlashcards)

	group.POST("/flashcards/review", fc.Review)

	group.GET("/flashcards", fc.GetAll)

	group.GET("/flashcards/:id", fc.GetByID)

	group.POST("/flashcards", fc.Create)

	group.PUT("/flashcards/:id", fc.Update)

	group.DELETE("/flashcards/:id", fc.Delete)

}
