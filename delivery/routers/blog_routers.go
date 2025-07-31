package routers

import (
	"github.com/abeni-al7/blog-platform/delivery/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(router *gin.Engine, blogController *controllers.BlogController) {
	router.POST("/blogs", blogController.CreateBlog)
}
