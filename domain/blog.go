package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string `gorm:"type:varchar(500)" json:"title"`
	Content   string `json:"content"`
	ViewCount int    `json:"view_count"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	UserID    int64  `json:"user_id"` // Foreign key column
	// User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // GORM relation
	CreatedAt time.Time `json:"created_at"` // auto set on insert
	UpdatedAt time.Time `json:"updated_at"` // auto set on update
}

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
}

type IBlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog) error
}
