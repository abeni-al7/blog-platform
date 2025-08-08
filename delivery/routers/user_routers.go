package routers

import (
	"os"

	"github.com/blog-platform/delivery/controllers"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/repositories"
	"github.com/blog-platform/usecases"
	"github.com/gin-gonic/gin"
)

// AuthRoutes registers authentication and user management routes on the provided Gin router group.
// It sets up endpoints for user registration, login, profile retrieval, and admin-only user promotion and demotion, applying appropriate authentication and authorization middleware.
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
	group.GET("/users/:id", ao.AccountOwnerMiddleware(), uc.GetProfile)
	adminRoutes := group.Group("/users")
	adminRoutes.Use(ao.AuthMiddleware(), ao.AdminMiddleware())
	{
		adminRoutes.PUT("/:id/promote", uc.Promote)
		adminRoutes.PUT("/:id/demote", uc.Demote)
	}
}
