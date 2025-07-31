package controllers

import (
	"net/http"

	"github.com/abeni-al7/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUsecase domain.IBlogUsecase
}

func NewBlogController(blogUsecase domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	var blog domain.Blog
	err := ctx.ShouldBindJSON(&blog)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.blogUsecase.CreateBlog(ctx.Request.Context(), &blog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog": blog})
}
