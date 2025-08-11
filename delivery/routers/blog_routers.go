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
	bc := controllers.NewBlogController(uu)

	blogRoutes := router.Group("/blogs")
	blogRoutes.Use(ao.AuthMiddleware())
	{
		blogRoutes.POST("", bc.CreateBlog)
		blogRoutes.GET("/:id", bc.GetBlogByID)
		blogRoutes.GET("", bc.GetBlogs)
		blogRoutes.DELETE("/:id", bc.DeleteBlog)
		blogRoutes.PATCH("/:id", bc.UpdateBlog)
		blogRoutes.GET("/paginated", bc.FetchPaginatedBlogs)
		blogRoutes.POST("/ideas", bc.GenerateBlogIdeas)
		blogRoutes.POST("/improve", bc.SuggestBlogImprovements)
		blogRoutes.GET("/filter", bc.FilterBlogs)
	}

}
