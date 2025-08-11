package repositories

import (
	"context"
	"errors"
	"fmt"
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
