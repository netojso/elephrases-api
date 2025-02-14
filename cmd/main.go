package main

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/config"
	_ "github.com/netojso/elephrases-api/docs"
	"github.com/netojso/elephrases-api/internal/adapters/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Elephrases API
// @version 1
// @description This is the Elephrases API documentation.
// @host localhost:8081
// @schemes http
// @produce json
// @consumes json

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	db := config.NewPostgresDatabase(cfg)

	defer config.ClosePostgresDBConnection(db)

	gin := gin.Default()

	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Setup(cfg, db, gin)

	gin.Run(cfg.ServerAddress)
}
