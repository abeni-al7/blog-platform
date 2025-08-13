package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	tokenInfra domain.IJWTInfrastructure
	blogRepo   domain.IBlogRepository
}

func NewMiddleware(tokenInfra domain.IJWTInfrastructure, blogRepo domain.IBlogRepository) *Middleware {
	return &Middleware{tokenInfra: tokenInfra, blogRepo: blogRepo}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		claims, err := m.tokenInfra.ValidateAccessToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		userID, _ := strconv.ParseInt(claims.UserID, 10, 64)
		role := claims.UserRole

		ctx.Set("user_id", userID)
		ctx.Set("role", role)
		ctx.Next()
	}
}

func (m *Middleware) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, ok := ctx.Get("role")
		if !ok || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to access this route"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *Middleware) AccountOwnerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		userID, ok := ctx.Get("user_id")

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil || !ok {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to access this route"})
			ctx.Abort()
			return
		}

		userIDInt, ok := userID.(int64)
		if !ok || userIDInt != idInt {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to access this route"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *Middleware) BlogAuthorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDVal, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		blogIDStr := ctx.Param("id")
		blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
		if err != nil || blogID <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
			ctx.Abort()
			return
		}

	authorID, err := m.blogRepo.GetBlogAuthorID(ctx.Request.Context(), blogID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			ctx.Abort()
			return
		}

		if authorID != userIDVal.(int64) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to access this route"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
