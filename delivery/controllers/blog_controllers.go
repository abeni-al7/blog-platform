package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUsecase domain.IBlogUsecase
}

func NewBlogController(uc domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase: uc}
}

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Tags    string `json:"tags" binding:"required"`
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int64)

	var req CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog := domain.Blog{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	// Split tags by comma and trim spaces
	var tags []string
	for _, tag := range strings.Split(req.Tags, ",") {
		t := strings.TrimSpace(tag)
		if t != "" {
			tags = append(tags, t)
		}
	}

	er := c.blogUsecase.CreateBlog(ctx.Request.Context(), &blog, tags)
	if er != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog": blog})
}

func (c *BlogController) GetBlogByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	blog, err := c.blogUsecase.FetchBlogByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	ctx.JSON(http.StatusOK, blog)
}

func (c *BlogController) GetBlogs(ctx *gin.Context) {
	blogs, err := c.blogUsecase.FetchAllBlogs(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blogs"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"blogs": blogs})

}

func (bc *BlogController) DeleteBlog(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	err = bc.blogUsecase.DeleteBlog(c.Request.Context(), blogID, userID)
	if err != nil {
		if err.Error() == "blog not found or not owned by user" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})
}

func (h *BlogController) FetchPaginatedBlogs(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	blogs, total, err := h.blogUsecase.FetchPaginatedBlogs(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch paginated blogs"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": blogs, "total": total, "page": page, "limit": limit, "total_pages": (total + int64(limit) - 1) / int64(limit)})
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

func (bc *BlogController) FilterBlogs(c *gin.Context) {
	// Parse query params
	title := c.Query("title")
	userIDStr := c.Query("user_id")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	var userIDPtr *int64
	if userIDStr != "" {
		uid, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userIDPtr = &uid
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	filter := domain.BlogFilter{
		TitleContains: title,
		UserID:        userIDPtr,
		Limit:         limit,
		Offset:        offset,
	}

	blogs, err := bc.blogUsecase.FetchBlogsByFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blogs)
}
