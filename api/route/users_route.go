package route

import (
	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/api/controller"
	"github.com/netojso/elephrases-api/repository"
	"github.com/netojso/elephrases-api/usecase"
	"gorm.io/gorm"
)

func NewUserRouter(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	uc := &controller.UserController{
		UserUseCase: usecase.NewUserUsecase(ur),
	}

	group.GET("/users", uc.Fetch)

	group.GET("/users/email/:email", uc.GetByEmail)

	group.GET("/users/:id", uc.GetUserByID)

	// update role
	group.PUT("/users/:id", uc.UpdateUser)

	group.DELETE("/users/:id", uc.DeleteUser)

}
