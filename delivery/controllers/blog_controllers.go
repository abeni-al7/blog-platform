package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

// helper: extract authenticated user ID from context
func getUserIDFromContext(ctx *gin.Context) (int64, bool) {
	// common keys set by auth middleware
	candidates := []string{"userID", "user_id", "uid"}
	for _, k := range candidates {
		if v, ok := ctx.Get(k); ok {
			switch t := v.(type) {
			case int64:
				return t, true
			case int:
				return int64(t), true
			case float64:
				return int64(t), true
			case string:
				if id, err := strconv.ParseInt(t, 10, 64); err == nil {
					return id, true
				}
			}
		}
	}
	return 0, false
}

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

func (c *BlogController) TrackView(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}
	if err := c.blogUsecase.TrackView(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to track view"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "view tracked"})
}

func (c *BlogController) LikeBlog(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	userID, ok := getUserIDFromContext(ctx)
	if !ok || userID <= 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.blogUsecase.LikeBlog(ctx.Request.Context(), id, userID); err != nil {
		// Map known errors to proper HTTP statuses.
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "not found"):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		case strings.Contains(msg, "already liked"):
			ctx.JSON(http.StatusConflict, gin.H{"error": "already liked"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like blog"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "liked"})
}

func (c *BlogController) UnlikeBlog(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	userID, ok := getUserIDFromContext(ctx)
	if !ok || userID <= 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.blogUsecase.UnlikeBlog(ctx.Request.Context(), id, userID); err != nil {
		// Mirror LikeBlog error mapping
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "not found"):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		case strings.Contains(msg, "not liked"):
			ctx.JSON(http.StatusConflict, gin.H{"error": "not liked"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike blog"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

func (c *BlogController) GetPopularity(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}
	views, likes, err := c.blogUsecase.GetPopularity(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get popularity"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"view_count": views, "likes": likes})
}
