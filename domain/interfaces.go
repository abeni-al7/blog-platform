package domain

import (
	"context"
)

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
	FindOrCreateTag(ctx context.Context, tagName string) (int64, error)
	LinkTagToBlog(ctx context.Context, blogID int64, tagID int64) error

	DeleteBlog(ctx context.Context, ID int64) error
}

type IBlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog, tags []string, userID int64) error
	DeleteBlog(ctx context.Context, ID int64, userID int64) error
}

type IJWTInfrastructure interface {
	GenerateAccessToken(userID string, userRole string) (string, error)
	GenerateRefreshToken(userID string, userRole string) (string, error)
	ValidateAccessToken(authHeader string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

type ITokenRepository interface {
	FetchByContent(content string) (Token, error)
}

type IPasswordInfrastructure interface {
	HashPassword(password string) (string, error)
	ComparePassword(correctPassword []byte, inputPassword []byte) error
}
