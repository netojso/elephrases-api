package router

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/config"
	"github.com/netojso/elephrases-api/internal/adapters/middleware"
	"gorm.io/gorm"
)

func Setup(env *config.Env, db *gorm.DB, gin *gin.Engine) {
	publicRouter := gin.Group("")
	// All Public APIs
	NewAuthRouter(env, db, publicRouter)

	protectedRouter := gin.Group("")
	// Middleware to verify AccessToken

	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	NewFlashcardRouter(env, db, protectedRouter)
	NewUsersRouter(db, protectedRouter)
	NewDeckRouter(db, protectedRouter)
}
