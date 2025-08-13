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
	ao := infrastructure.NewMiddleware(js, ur)
	ai := infrastructure.NewChatGPTAIService()
	uu := usecases.NewBlogUsecase(ur, ai)
	bc := controllers.NewBlogController(uu)

	blogRoutes := router.Group("/blogs")
	blogRoutes.Use(ao.AuthMiddleware())
	{
		blogRoutes.POST("", bc.CreateBlog)
		blogRoutes.GET("/:id", bc.GetBlogByID)
		blogRoutes.GET("", bc.GetBlogs)
		blogRoutes.DELETE("/:id", ao.BlogAuthorMiddleware(), bc.DeleteBlog)
		blogRoutes.PATCH("/:id", ao.BlogAuthorMiddleware(), bc.UpdateBlog)
		blogRoutes.GET("/paginated", bc.FetchPaginatedBlogs)
		blogRoutes.GET("/search", bc.SearchBlogs)
		blogRoutes.POST("/:id/view", bc.TrackView)
		blogRoutes.POST("/:id/like", bc.LikeBlog)
		blogRoutes.DELETE("/:id/like", bc.UnlikeBlog)
		blogRoutes.GET("/:id/popularity", bc.GetPopularity)
		blogRoutes.POST("/ideas", bc.GenerateBlogIdeas)
		blogRoutes.POST("/improve", bc.SuggestBlogImprovements)
		blogRoutes.GET("/filter", bc.FilterBlogs)
		// comments
		blogRoutes.POST("/:id/comments", bc.AddComment)
		blogRoutes.GET("/:id/comments", bc.ListComments)
	}

}
