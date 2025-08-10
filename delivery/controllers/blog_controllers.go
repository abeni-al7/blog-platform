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
	if err := c.blogUsecase.LikeBlog(ctx.Request.Context(), id, 0); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like blog"})
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
	if err := c.blogUsecase.UnlikeBlog(ctx.Request.Context(), id, 0); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike blog"})
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

func (h *BlogController) SearchBlogs(ctx *gin.Context) {
	q := ctx.Query("q")
	if strings.TrimSpace(q) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	blogs, total, err := h.blogUsecase.SearchBlogs(ctx.Request.Context(), q, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search blogs"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"blogs": blogs,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
