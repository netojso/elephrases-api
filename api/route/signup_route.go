package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/api/controller"
	"github.com/netojso/elephrases-api/bootstrap"
	"github.com/netojso/elephrases-api/repository"
	"github.com/netojso/elephrases-api/usecase"
	"gorm.io/gorm"
)

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur),
		Env:           env,
	}
	group.POST("/signup", sc.Signup)
}
