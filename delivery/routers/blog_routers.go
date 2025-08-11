package routers

import (
	"github.com/blog-platform/delivery/controllers"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/repositories"
	"github.com/blog-platform/usecases"
	"github.com/gin-gonic/gin"
)

func BlogRoutes(group *gin.RouterGroup) {
	DB := repositories.DB
	blogRepo := repositories.NewBlogRepository(DB)
	aiService := infrastructure.NewChatGPTAIService()
	blogUsecase := usecases.NewBlogUsecase(blogRepo, aiService)
	blogController := controllers.NewBlogController(blogUsecase)

	// Grouped under /blogs
	blogs := group.Group("/blogs")
	{
		blogs.POST("", blogController.CreateBlog)
		blogs.GET("/:id", blogController.GetBlogByID)
		blogs.GET("", blogController.GetBlogs)
		blogs.DELETE("/:id", blogController.DeleteBlog)
		blogs.GET("/paginated", blogController.FetchPaginatedBlogs)
		blogs.POST("/ideas", blogController.GenerateBlogIdeas)
		blogs.POST("/improve", blogController.SuggestBlogImprovements)
		blogs.GET("/filter", blogController.FilterBlogs)
	}
}
