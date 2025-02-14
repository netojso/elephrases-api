package router

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/handler"
	repository "github.com/netojso/elephrases-api/internal/adapters/repository/deck"
	"github.com/netojso/elephrases-api/internal/core/service"
	"gorm.io/gorm"
)

func NewDeckRouter(db *gorm.DB, r *gin.RouterGroup) {
	repository := repository.NewDeckRepository(db)
	service := service.NewDeckUsecase(repository)
	handler := handler.NewDeckHandler(service)

	r.GET("/decks", handler.GetAll)
	r.GET("/decks/:id", handler.GetByID)
	r.POST("/decks", handler.Create)
	r.PUT("/decks/:id", handler.Update)
	r.DELETE("/decks/:id", handler.Delete)

}
