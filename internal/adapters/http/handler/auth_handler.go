package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/dto"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	service portservice.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(service portservice.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// @Summary Login a user
// @Description Login a user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body dto.LoginDTO true "Login credentials"
// @Success 200 {object} dto.AuthResponseDTO
// @Router /login [post]
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

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body dto.RegisterDTO true "Registration details"
// @Success 200 {object} dto.AuthResponseDTO
// @Router /register [post]
func (ah AuthHandler) Register(ctx *gin.Context) {
	var body dto.RegisterDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session, err := ah.service.Register(body.Email, body.Password)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response := session.ToMap()

	ctx.JSON(http.StatusOK, response)
}
