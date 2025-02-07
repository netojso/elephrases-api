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

	router := r.Group("/decks")

	router.GET("/", handler.GetAll)
	router.GET("/:id", handler.GetByID)
	router.POST("/", handler.Create)
	router.PUT("/:id", handler.Update)
	router.DELETE("/:id", handler.Delete)

}
