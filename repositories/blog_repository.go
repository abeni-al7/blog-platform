package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/infrastructure"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
	c  *infrastructure.Cache
}

// internal cached value for paginated results
type pagedBlogs struct {
	Blogs []*domain.Blog
	Total int64
}

func NewBlogRepository(db *gorm.DB) domain.IBlogRepository {
	return &BlogRepository{db: db, c: infrastructure.NewCache()}
}

func (r *BlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	if err := r.db.WithContext(ctx).Create(blog).Error; err != nil {
		return err
	}
	// invalidate caches
	r.c.Clear()
	return nil
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
	if err := r.db.WithContext(ctx).Create(&tagBlog).Error; err != nil {
		return err
	}
	r.c.Clear()
	return nil
}
func (r *BlogRepository) FetchByID(ctx context.Context, id int64) (*domain.Blog, error) {
	key := fmt.Sprintf("blog:%d", id)
	if v, ok := r.c.Get(key); ok {
		if b, ok2 := v.(*domain.Blog); ok2 {
			return b, nil
		}
	}
	var blog domain.Blog
	if err := r.db.WithContext(ctx).Preload("User").Preload("Tags").First(&blog, id).Error; err != nil {
		return nil, err
	}
	r.c.Set(key, &blog, 5*time.Minute)
	return &blog, nil
}

func (r *BlogRepository) FetchAll(ctx context.Context) ([]*domain.Blog, error) {
	const key = "blogs:all"
	if v, ok := r.c.Get(key); ok {
		if bs, ok2 := v.([]*domain.Blog); ok2 {
			return bs, nil
		}
	}
	var blogs []*domain.Blog
	if err := r.db.WithContext(ctx).Preload("User").Preload("Tags").Find(&blogs).Error; err != nil {
		return nil, err
	}
	r.c.Set(key, blogs, 2*time.Minute)
	return blogs, nil
}

func (r *BlogRepository) DeleteByID(ctx context.Context, ID int64, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", ID, userID).
		Delete(&domain.Blog{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("blog not found")
	}
	r.c.Clear()
	return result.Error
}

func (r *BlogRepository) UpdateByID(ctx context.Context, id int64, userID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).Model(&domain.Blog{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("blog not found")
	}
	// invalidate caches after successful update
	r.c.Clear()
	return nil
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(fb *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return fb.Offset(offset).Limit(limit)
	}
}
func (r *BlogRepository) FetchPaginatedBlogs(ctx context.Context, page, limit int) ([]*domain.Blog, int64, error) {
	key := fmt.Sprintf("blogs:p=%d:l=%d", page, limit)
	if v, ok := r.c.Get(key); ok {
		if pb, ok2 := v.(*pagedBlogs); ok2 {
			return pb.Blogs, pb.Total, nil
		}
	}
	var blogs []*domain.Blog
	var total int64

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return r.db.WithContext(gctx).Model(&domain.Blog{}).Count(&total).Error
	})
	g.Go(func() error {
		return r.db.WithContext(gctx).
			Preload("User").
			Preload("Tags").
			Scopes(Paginate(page, limit)).
			Find(&blogs).Error
	})
	if err := g.Wait(); err != nil {
		return nil, 0, err
	}
	r.c.Set(key, &pagedBlogs{Blogs: blogs, Total: total}, 1*time.Minute)
	return blogs, total, nil

}

func (r *BlogRepository) FetchByFilter(ctx context.Context, filter domain.BlogFilter) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	query := r.db.WithContext(ctx).Model(&domain.Blog{})

	if filter.TitleContains != "" {
		query = query.Where("title ILIKE ?", "%"+filter.TitleContains+"%")
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset >= 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&blogs).Error
	if err != nil {
		return nil, err
	}

	return blogs, nil
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
	// Portable clamp to zero
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var b domain.Blog
		if err := tx.Select("id, likes").First(&b, blogID).Error; err != nil {
			return err
		}
		if b.Likes > 0 {
			b.Likes--
		}
		return tx.Model(&domain.Blog{}).Where("id = ?", blogID).Update("likes", b.Likes).Error
	})
}

func (r *BlogRepository) GetPopularity(ctx context.Context, blogID int64) (int, int, error) {
	var b domain.Blog
	if err := r.db.WithContext(ctx).Select("id, view_count, likes").First(&b, blogID).Error; err != nil {
		return 0, 0, err
	}
	return b.ViewCount, b.Likes, nil
}

func (r *BlogRepository) SearchBlogs(ctx context.Context, query string, page, limit int) ([]*domain.Blog, int64, error) {
	var (
		blogs []*domain.Blog
		total int64
	)

	q := strings.TrimSpace(query)
	if q == "" {
		return []*domain.Blog{}, 0, nil
	}
	pattern := "%" + strings.ToLower(q) + "%"

	db := r.db.WithContext(ctx).Model(&domain.Blog{}).
		Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ?", pattern, pattern)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}
	// r.c.Set(key, &pagedBlogs{Blogs: blogs, Total: total}, 1*time.Minute)
	return blogs, total, nil

}

func (r *BlogRepository) AddComment(ctx context.Context, blogID, userID int64, content string) (*domain.Comment, error) {
	c := &domain.Comment{
		BlogID:  blogID,
		UserID:  userID,
		Content: content,
	}

	// Create the comment
	if err := r.db.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Blog").
		Preload("Blog.User").
		First(c, c.ID).Error; err != nil {
		return nil, err
	}

	return c, nil
}

func (r *BlogRepository) ListComments(ctx context.Context, blogID int64, page, limit int) ([]*domain.Comment, int64, error) {
	var (
		comments []*domain.Comment
		total    int64
	)

	q := r.db.WithContext(ctx).Model(&domain.Comment{}).Where("blog_id = ?", blogID)

	// Count total
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Fetch with preloaded associations
	if err := q.
		Preload("User").
		Preload("Blog").
		Preload("Blog.User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
