package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/dto"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
)

type AuthHandler struct {
	service portservice.AuthService
}

func NewAuthHandler(service portservice.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (ah AuthHandler) Login(ctx *gin.Context) {
	var body dto.LoginDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session, err := ah.service.Login(body.Email, body.Password)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response := session.ToMap()

	ctx.JSON(http.StatusOK, response)
}

func (ah AuthHandler) Register(ctx *gin.Context) {
	var body dto.RegisterDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := ah.service.Register(body.Email, body.Password)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response := user.ToMap()

	ctx.JSON(http.StatusOK, response)
}
