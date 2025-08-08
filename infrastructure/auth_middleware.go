package infrastructure

import (
	"net/http"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	tokenInfra domain.IJWTInfrastructure
}

func NewMiddleware(tokenInfra domain.IJWTInfrastructure) *Middleware {
	return &Middleware{tokenInfra: tokenInfra}
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

		userID := claims.UserID
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

		if !ok || userID != id {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to access this route"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
