package main

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/config"
	"github.com/netojso/elephrases-api/internal/adapters/router"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	db := config.NewPostgresDatabase(cfg)

	defer config.ClosePostgresDBConnection(db)

	gin := gin.Default()

	router.Setup(cfg, db, gin)

	gin.Run(cfg.ServerAddress)
}
