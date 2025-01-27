package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netojso/go-api-template/domain"
)

type UserController struct {
	UserUseCase domain.UserUsecase
}

func (uc *UserController) Fetch(ctx *gin.Context) {
	users, err := uc.UserUseCase.Fetch()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	user, err := uc.UserUseCase.GetByEmail(email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given email"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := uc.UserUseCase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given id"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) CreateUser(ctx *gin.Context, user domain.User) error {
	return uc.UserUseCase.Create(user)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedUser domain.User

	if err := ctx.ShouldBind(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	_, err := uc.UserUseCase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given id"})
		return
	}

	uc.UserUseCase.UpdateUser(id, updatedUser)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uc.UserUseCase.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given id"})
		return
	}

	uc.UserUseCase.DeleteUser(id)
}
