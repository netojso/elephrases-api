package router

import (
	"time"

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
	NewUsersRouter(env, db, protectedRouter)
	NewDeckRouter(env, db, protectedRouter)
}

// timeout converte CONTEXT_TIMEOUT (segundos) no timeout aplicado a cada request.
func timeout(env *config.Env) time.Duration {
	return time.Duration(env.ContextTimeout) * time.Second
}
