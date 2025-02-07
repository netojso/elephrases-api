package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg"
)

type UserHandler struct {
	service portservice.UserService
}

func NewUserHandler(service portservice.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (uc *UserHandler) Fetch(ctx *gin.Context) {
	users, err := uc.service.Fetch()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (uc *UserHandler) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	user, err := uc.service.GetByEmail(email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User not found with the given email"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserHandler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := uc.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User not found with the given id"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserHandler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := uc.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User not found with the given id"})
		return
	}

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if err := uc.service.UpdateUser(id, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
}

func (uc *UserHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uc.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User not found with the given id"})
		return
	}

	uc.service.DeleteUser(id)
}
