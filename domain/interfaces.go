package domain

import (
	"context"
)

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
	FindOrCreateTag(ctx context.Context, tagName string) (int64, error)
	LinkTagToBlog(ctx context.Context, blogID int64, tagID int64) error
	FetchAll(ctx context.Context) ([]*Blog, error)
}

type IAIService interface {
	GenerateBlogIdeas(topic string) (string, error)
	SuggestBlogImprovements(content string) (string, error)
}

type IBlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog, tags []string) error
	FetchAllBlogs(ctx context.Context) ([]*Blog, error)
	GenerateBlogIdeas(topic string) (string, error)
	SuggestBlogImprovements(content string) (string, error)
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

type IEmailInfrastructure interface {
	SendEmail(to []string, subject string, body string) error
}

type IUserUsecase interface {
	Register(user *User) (User, error)
	ActivateAccount(id string) error
}

type IUserRepository interface {
	Register(user *User) (User, error)
	FetchByUsername(username string) (User, error)
	FetchByEmail(email string) (User, error)
	ActivateAccount(idStr string) error
	Fetch(idStr string) (User, error)
}

type IUserController interface {
	Register(ctx *context.Context)
	ActivateAccount(ctx *context.Context)
}
