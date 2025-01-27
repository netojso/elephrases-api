package route

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/api/controller"
	"github.com/netojso/elephrases-api/repository"
	"github.com/netojso/elephrases-api/usecase"
	"gorm.io/gorm"
)

func NewDeckRoute(db *gorm.DB, group *gin.RouterGroup) {
	fr := repository.NewDeckRepository(db)
	fc := &controller.DeckController{
		DeckUsecase: usecase.NewDeckUsecase(fr),
	}

	group.GET("/decks", fc.GetAll)

	group.GET("/decks/:id", fc.GetByID)

	group.POST("/decks", fc.Create)

	group.PUT("/decks/:id", fc.Update)

	group.DELETE("/decks/:id", fc.Delete)

}
