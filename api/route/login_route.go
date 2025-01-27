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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	lc := &controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur),
		Env:          env,
	}
	group.POST("/login", lc.Login)
}
