package domain

import (
	"context"
)

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
	FindOrCreateTag(ctx context.Context, tagName string) (int64, error)
	LinkTagToBlog(ctx context.Context, blogID int64, tagID int64) error
	FetchByID(ctx context.Context, id int64) (*Blog, error)
	FetchAll(ctx context.Context) ([]*Blog, error)
	IncrementView(ctx context.Context, blogID int64) error
	AddLike(ctx context.Context, blogID int64, userID int64) error // userID optional, ignored
	RemoveLike(ctx context.Context, blogID int64, userID int64) error
	GetPopularity(ctx context.Context, blogID int64) (views int, likes int, err error)

	FetchPaginatedBlogs(ctx context.Context, page int, limit int) ([]*Blog, int64, error)
=========
	DeleteByID(ctx context.Context, ID int64, userID string) error
	UpdateByID(ctx context.Context, id int64, userID string, updates map[string]interface{}) error
	FetchByFilter(ctx context.Context, filter BlogFilter) ([]*Blog, error)
}

type IAIService interface {
	GenerateBlogIdeas(topic string) (string, error)
	SuggestBlogImprovements(content string) (string, error)
}

type IBlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog, tags []string) error
	FetchBlogByID(ctx context.Context, id int64) (*Blog, error)
	FetchAllBlogs(ctx context.Context) ([]*Blog, error)
	DeleteBlog(ctx context.Context, ID int64, userID string) error
	FetchPaginatedBlogs(ctx context.Context, page int, limit int) ([]*Blog, int64, error)
	TrackView(ctx context.Context, blogID int64) error
	LikeBlog(ctx context.Context, blogID, userID int64) error
	UnlikeBlog(ctx context.Context, blogID, userID int64) error
	GetPopularity(ctx context.Context, blogID int64) (views int, likes int, err error)
	SearchBlogs(ctx context.Context, query string, page, limit int) ([]*Blog, int64, error)
	GenerateBlogIdeas(topic string) (string, error)
	SuggestBlogImprovements(content string) (string, error)
	UpdateBlog(ctx context.Context, id int64, userID string, updates map[string]interface{}) error
	FetchBlogsByFilter(ctx context.Context, filter BlogFilter) ([]*Blog, error)
}

type IJWTInfrastructure interface {
	GenerateAccessToken(userID string, userRole string) (string, error)
	GenerateRefreshToken(userID string, userRole string) (string, error)
	ValidateAccessToken(authHeader string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

type ITokenRepository interface {
	FetchByContent(content string) (Token, error)
	Save(token *Token) error
	Delete(content string) error
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
	Login(identifier string, password string) (string, string, error)
	GetUserProfile(userID int64) (*User, error)
	Promote(id string) error
	Demote(id string) error
	UpdateUserProfile(userID int64, updates map[string]interface{}) error
	RefreshToken(authHeader string) (string, string, error)
	ResetPassword(userID string, oldPassword string, newPassword string) error
	ForgotPassword(email string) error
	UpdatePasswordDirect(userID string, newPassword string, token string) error
	Logout(authHeader string) error
}

type IUserRepository interface {
	Register(user *User) (User, error)
	FetchByUsername(username string) (User, error)
	FetchByEmail(email string) (User, error)
	ActivateAccount(idStr string) error
	Fetch(idStr string) (User, error)
	GetUserProfile(userID int64) (*User, error)
	Promote(idStr string) error
	Demote(idStr string) error
	UpdateUserProfile(userID int64, updates map[string]interface{}) error
	ResetPassword(idStr string, newPassword string) error
	CountUsers() (int64, error)
}

type IUserController interface {
	Register(ctx *context.Context)
	ActivateAccount(ctx *context.Context)
	Login(ctx *context.Context)

	GetProfile(ctx *context.Context)
	UpdateProfile(ctx *context.Context)
	RefreshToken(ctx *context.Context)
	ResetPassword(ctx *context.Context)

	ForgotPassword(ctx *context.Context)
	UpdatePasswordDirect(ctx *context.Context)
	Logout(ctx *context.Context)
}
