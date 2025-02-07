package router

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/config"
	"github.com/netojso/elephrases-api/internal/adapters/http/handler"
	repository "github.com/netojso/elephrases-api/internal/adapters/repository/auth"
	"github.com/netojso/elephrases-api/internal/core/service"
	"gorm.io/gorm"
)

func NewAuthRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	repo := repository.NewAuthRepository(db)
	auth_service := service.NewAuthService(repo, env)
	auth_handler := handler.NewAuthHandler(auth_service)

	group.POST("/login", auth_handler.Login)
	group.POST("/register", auth_handler.Register)
}
