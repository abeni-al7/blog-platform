package routers

import (
	"os"

	"github.com/blog-platform/delivery/controllers"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/repositories"
	"github.com/blog-platform/usecases"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(group *gin.RouterGroup) {
	DB := repositories.DB
	ur := repositories.NewUserRepository(DB)
	ei := infrastructure.NewSMTPEmailService()
	pi := infrastructure.NewPasswordInfrastructure()
	tr := repositories.NewTokenRepository(DB)
	js := infrastructure.NewJWTInfrastructure([]byte(os.Getenv("JWT_ACCESS_SECRET")), []byte(os.Getenv("JWT_REFRESH_SECRET")), tr)
	uu := usecases.NewUserUsecase(ur, ei, pi, js, tr)
	uc := controllers.NewUserController(uu)
	ao := infrastructure.NewMiddleware(js)

	group.POST("/register", uc.Register)
	group.POST("/login", uc.Login)
	group.POST("/token/refresh", uc.RefreshToken)
	group.POST("/reset-password", ao.AuthMiddleware(), uc.ResetPassword)
	group.POST("/forgot-password", uc.ForgotPassword)
	group.POST("/password/:id/update", uc.UpdatePasswordDirect)
	group.GET("/users/:id", ao.AccountOwnerMiddleware(), uc.GetProfile)
	group.PATCH("/users/:id", ao.AccountOwnerMiddleware(), uc.UpdateProfile)
}
