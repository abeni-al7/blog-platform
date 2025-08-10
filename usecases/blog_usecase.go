package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/blog-platform/domain"
)

type blogUsecase struct {
	blogRepo domain.IBlogRepository
}

func NewBlogUsecase(repo domain.IBlogRepository) domain.IBlogUsecase {
	return &blogUsecase{
		blogRepo: repo,
	}
}

func (uc blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, tags []string) error {
	// prevent empty strings from being added
	if blog.Title == "" || blog.Content == "" {
		return errors.New("title and content cannot be empty")
	}

	err := uc.blogRepo.Create(ctx, blog)

	if err != nil {
		return errors.New("failed to create blog")
	}
	// Ensure blog ID is populated after creation
	if blog.ID == 0 {
		return errors.New("blog ID not set after creation")
	}

	for _, tag := range tags {
		if tag == "" {
			continue // skip empty tags
		}

		tagID, err := uc.blogRepo.FindOrCreateTag(ctx, tag)
		if err != nil {
			return fmt.Errorf("failed to find or create tag '%s': %w", tag, err)
		}

		err = uc.blogRepo.LinkTagToBlog(ctx, int64(blog.ID), tagID)
		if err != nil {
			return fmt.Errorf("failed to link tag '%s' to blog: %w", tag, err)
		}
	}

	return nil
}

func (uc blogUsecase) FetchBlogByID(ctx context.Context, id int64) (*domain.Blog, error) {
	if id <= 0 {
		return nil, errors.New("invalid blog ID")
	}

	blog, err := uc.blogRepo.FetchByID(ctx, id)
	if err != nil {
		return nil, errors.New("failed to fetch blog")
	}

	return blog, nil
}

func (uc *blogUsecase) FetchAllBlogs(ctx context.Context) ([]*domain.Blog, error) {
	blogs, err := uc.blogRepo.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blogs: %w", err)
	}
	return blogs, nil
}

func (uc *blogUsecase) FetchPaginatedBlogs(ctx context.Context, page, limit int) ([]*domain.Blog, int64, error) {
	return uc.blogRepo.FetchPaginatedBlogs(ctx, page, limit)
}

func (uc *blogUsecase) TrackView(ctx context.Context, blogID int64) error {
	if blogID <= 0 {
		return errors.New("invalid blog ID")
	}
	return uc.blogRepo.IncrementView(ctx, blogID)
}

func (uc *blogUsecase) LikeBlog(ctx context.Context, blogID, userID int64) error {
	if blogID <= 0 {
		return errors.New("invalid blog ID")
	}
	return uc.blogRepo.AddLike(ctx, blogID, userID)
}

func (uc *blogUsecase) UnlikeBlog(ctx context.Context, blogID, userID int64) error {
	if blogID <= 0 {
		return errors.New("invalid blog ID")
	}
	return uc.blogRepo.RemoveLike(ctx, blogID, userID)
}

func (uc *blogUsecase) GetPopularity(ctx context.Context, blogID int64) (int, int, error) {
	if blogID <= 0 {
		return 0, 0, errors.New("invalid blog ID")
	}
	return uc.blogRepo.GetPopularity(ctx, blogID)
}
