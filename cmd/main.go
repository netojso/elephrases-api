package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/netojso/go-api-template/api/route"
	"github.com/netojso/go-api-template/bootstrap"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Db
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, db, gin)

	gin.Run(env.ServerAddress)
}
