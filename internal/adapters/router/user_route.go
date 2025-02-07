package router

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/handler"
	repository "github.com/netojso/elephrases-api/internal/adapters/repository/users"
	"github.com/netojso/elephrases-api/internal/core/service"
	"gorm.io/gorm"
)

func NewUsersRouter(db *gorm.DB, group *gin.RouterGroup) {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	group.GET("/users", handler.Fetch)

	group.GET("/users/email/:email", handler.GetByEmail)

	group.GET("/users/:id", handler.GetUserByID)

	// update role
	group.PUT("/users/:id", handler.UpdateUser)

	group.DELETE("/users/:id", handler.DeleteUser)

}
