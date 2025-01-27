package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/netojso/go-api-template/api/controller"
	"github.com/netojso/go-api-template/bootstrap"
	"github.com/netojso/go-api-template/repository"
	"github.com/netojso/go-api-template/usecase"
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
