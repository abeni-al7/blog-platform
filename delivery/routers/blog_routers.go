package routers

import (
	"os"

	"github.com/blog-platform/delivery/controllers"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/repositories"
	"github.com/blog-platform/usecases"
	"github.com/gin-gonic/gin"
)

func BlogRoutes(router *gin.RouterGroup) {

	DB := repositories.DB
	ur := repositories.NewBlogRepository(DB)
	tr := repositories.NewTokenRepository(DB)
	js := infrastructure.NewJWTInfrastructure([]byte(os.Getenv("JWT_ACCESS_SECRET")), []byte(os.Getenv("JWT_REFRESH_SECRET")), tr)
	ao := infrastructure.NewMiddleware(js)
	ai := infrastructure.NewChatGPTAIService()
	uu := usecases.NewBlogUsecase(ur, ai)
	uc := controllers.NewBlogController(uu)

	blogRoutes := router.Group("/blogs")
	blogRoutes.Use(ao.AuthMiddleware())
	{
		blogRoutes.POST("", uc.CreateBlog)
		blogRoutes.GET("/:id", uc.GetBlogByID)
		blogRoutes.GET("", uc.GetBlogs)
		blogRoutes.DELETE("/:id", uc.DeleteBlog)
		blogRoutes.PATCH("/:id", uc.UpdateBlog)
		blogRoutes.GET("/paginated", uc.FetchPaginatedBlogs)
		blogRoutes.POST("/ideas", uc.GenerateBlogIdeas)
		blogRoutes.POST("/improve", uc.SuggestBlogImprovements)
	}

}
