package routers

import (
	"github.com/blog-platform/delivery/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(router *gin.Engine, blogController *controllers.BlogController) {
	router.POST("/blogs", blogController.CreateBlog)

	router.GET("/blogs/:id", blogController.GetBlogByID)

	router.GET("/blogs", blogController.GetBlogs)

	router.GET("/blogs/paginated", blogController.FetchPaginatedBlogs)
	router.GET("/blogs/search", blogController.SearchBlogs)
	router.POST("/blogs/:id/view", blogController.TrackView)
	router.POST("/blogs/:id/like", blogController.LikeBlog)
	router.DELETE("/blogs/:id/like", blogController.UnlikeBlog)
	router.GET("/blogs/:id/popularity", blogController.GetPopularity)
}
