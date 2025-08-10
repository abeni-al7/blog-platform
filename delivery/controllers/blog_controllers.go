package controllers

import (
	"net/http"
	"strings"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUsecase domain.IBlogUsecase
}

func NewBlogController(blogUsecase domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase}
}

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Tags    string `json:"tags" binding:"required"`
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	var req CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog := domain.Blog{
		Title:   req.Title,
		Content: req.Content,
	}

	// Split tags by comma and trim spaces
	var tags []string
	for _, tag := range strings.Split(req.Tags, ",") {
		t := strings.TrimSpace(tag)
		if t != "" {
			tags = append(tags, t)
		}
	}

	err := c.blogUsecase.CreateBlog(ctx.Request.Context(), &blog, tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog": blog})
}

func (c *BlogController) GetBlogs(ctx *gin.Context) {
	blogs, err := c.blogUsecase.FetchAllBlogs(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blogs"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

type BlogIdeaRequest struct {
	Topic string `json:"topic" binding:"required"`
}

func (c *BlogController) GenerateBlogIdeas(ctx *gin.Context) {
	var req BlogIdeaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ideas, err := c.blogUsecase.GenerateBlogIdeas(req.Topic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ideas": ideas})
}

type BlogImproveRequest struct {
	Content string `json:"content" binding:"required"`
}

func (c *BlogController) SuggestBlogImprovements(ctx *gin.Context) {
	var req BlogImproveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suggestion, err := c.blogUsecase.SuggestBlogImprovements(req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"suggestion": suggestion})
}
