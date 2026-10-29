package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/internal/adapters/http/dto"
	portservice "github.com/netojso/elephrases-api/internal/core/ports/service"
	"github.com/netojso/elephrases-api/pkg"
)

type UserHandler struct {
	timeout time.Duration
	service portservice.UserService
}

func NewUserHandler(service portservice.UserService, timeout time.Duration) *UserHandler {
	return &UserHandler{service: service, timeout: timeout}
}

// Fetch godoc
// @Summary Fetch all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.User
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users [get]
func (uc *UserHandler) Fetch(ctx *gin.Context) {
	reqCtx, cancel := requestContext(ctx, uc.timeout)
	defer cancel()

	users, err := uc.service.Fetch(reqCtx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetByEmail godoc
// @Summary Get user by email
// @Description Get a user by email
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email"
// @Success 200 {object} domain.User
// @Failure 404 {object} pkg.ErrorResponse
// @Router /users/email/{email} [get]
func (uc *UserHandler) GetByEmail(ctx *gin.Context) {
	reqCtx, cancel := requestContext(ctx, uc.timeout)
	defer cancel()

	email := ctx.Param("email")

	user, err := uc.service.GetByEmail(reqCtx, email)

	if err != nil {
		ctx.JSON(errorStatus(err), pkg.ErrorResponse{Message: userErrorMessage(err, "User not found with the given email")})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} domain.User
// @Failure 404 {object} pkg.ErrorResponse
// @Router /users/{id} [get]
func (uc *UserHandler) GetUserByID(ctx *gin.Context) {
	reqCtx, cancel := requestContext(ctx, uc.timeout)
	defer cancel()

	id := ctx.Param("id")

	user, err := uc.service.GetByID(reqCtx, id)

	if err != nil {
		ctx.JSON(errorStatus(err), pkg.ErrorResponse{Message: userErrorMessage(err, "User not found with the given id")})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Param user body dto.UpdateUserDTO true "User"
// @Success 200 {object} domain.User
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id} [put]
func (uc *UserHandler) UpdateUser(ctx *gin.Context) {
	reqCtx, cancel := requestContext(ctx, uc.timeout)
	defer cancel()

	id := ctx.Param("id")

	user, err := uc.service.GetByID(reqCtx, id)

	if err != nil {
		ctx.JSON(errorStatus(err), pkg.ErrorResponse{Message: userErrorMessage(err, "User not found with the given id")})
		return
	}

	var body dto.UpdateUserDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	user.FullName = body.FullName
	user.PhoneNumber = body.PhoneNumber

	if err := uc.service.UpdateUser(reqCtx, id, user); err != nil {
		ctx.JSON(errorStatus(err), pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 204
// @Failure 404 {object} pkg.ErrorResponse
// @Router /users/{id} [delete]
func (uc *UserHandler) DeleteUser(ctx *gin.Context) {
	reqCtx, cancel := requestContext(ctx, uc.timeout)
	defer cancel()

	id := ctx.Param("id")

	_, err := uc.service.GetByID(reqCtx, id)

	if err != nil {
		ctx.JSON(errorStatus(err), pkg.ErrorResponse{Message: userErrorMessage(err, "User not found with the given id")})
		return
	}

	if err := uc.service.DeleteUser(reqCtx, id); err != nil {
		ctx.JSON(errorStatus(err), pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// userErrorMessage mantém a mensagem amigável para 404 e expõe o erro nos demais casos.
func userErrorMessage(err error, notFound string) string {
	if errors.Is(err, pkg.ErrNotFound) {
		return notFound
	}
	return err.Error()
}
