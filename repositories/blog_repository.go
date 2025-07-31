package repositories

import (
	"context"

	"github.com/abeni-al7/blog-platform/domain"
	"gorm.io/gorm"
)

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) domain.IBlogRepository {
	return &blogRepository{db}
}

func (r *blogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}
