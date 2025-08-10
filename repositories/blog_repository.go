package repositories

import (
	"context"
	"errors"

	"github.com/blog-platform/domain"
	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) domain.IBlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *BlogRepository) FindOrCreateTag(ctx context.Context, tagName string) (int64, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).Where("name = ?", tagName).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tag = domain.Tag{Name: tagName}
		if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
			return 0, err
		}
		return tag.ID, nil
	} else if err != nil {
		return 0, err
	}
	return tag.ID, nil
}

func (r *BlogRepository) LinkTagToBlog(ctx context.Context, blogID int64, tagID int64) error {
	tagBlog := domain.Tag_Blog{
		BlogID: blogID,
		TagID:  tagID,
	}
	return r.db.WithContext(ctx).Create(&tagBlog).Error
}

func (r *BlogRepository) FetchByID(ctx context.Context, id int64) (*domain.Blog, error) {
	var blog domain.Blog
	if err := r.db.WithContext(ctx).Preload("User").Preload("Tags").First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) FetchAll(ctx context.Context) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	if err := r.db.WithContext(ctx).Preload("User").Preload("Tags").
		Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(fb *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return fb.Offset(offset).Limit(limit)
	}
}
func (r *BlogRepository) FetchPaginatedBlogs(ctx context.Context, page, limit int) ([]*domain.Blog, int64, error) {
	var blogs []*domain.Blog
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Tags").
		Scopes(Paginate(page, limit)).
		Find(&blogs).Error
	if err != nil {
		return nil, 0, err
	}
	return blogs, total, nil
}

func (r *BlogRepository) IncrementView(ctx context.Context, blogID int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Blog{}).
		Where("id = ?", blogID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *BlogRepository) AddLike(ctx context.Context, blogID int64, _ int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Blog{}).
		Where("id = ?", blogID).
		UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}

func (r *BlogRepository) RemoveLike(ctx context.Context, blogID int64, _ int64) error {
	tx := r.db.WithContext(ctx).
		Model(&domain.Blog{}).
		Where("id = ? AND likes > 0", blogID).
		UpdateColumn("likes", gorm.Expr("likes - 1"))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 1 {
		return nil
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Blog{}).
		Where("id = ?", blogID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return errors.New("not liked")
}

func (r *BlogRepository) GetPopularity(ctx context.Context, blogID int64) (int, int, error) {
	var b domain.Blog
	if err := r.db.WithContext(ctx).Select("id, view_count, likes").First(&b, blogID).Error; err != nil {
		return 0, 0, err
	}
	return b.ViewCount, b.Likes, nil
}
