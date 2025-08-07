package routers

import (
	"github.com/blog-platform/delivery/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(router *gin.Engine, blogController *controllers.BlogController) {
	router.POST("/blogs", blogController.CreateBlog)
	router.GET("/blogs", blogController.GetBlogs)
	router.POST("/blogs/ideas", blogController.GenerateBlogIdeas)
}
