package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/blog-platform/domain"
)

type blogUsecase struct {
	blogRepo  domain.IBlogRepository
	aiService domain.IAIService
}

func NewBlogUsecase(repo domain.IBlogRepository, aiService domain.IAIService) domain.IBlogUsecase {
	return &blogUsecase{
		blogRepo:  repo,
		aiService: aiService,
	}
}

func (uc blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, tags []string) error {
	// prevent empty strings from being added

	if blog.Title == "" || blog.Content == "" {
		return errors.New("title and content cannot be empty")
	}
	if blog.UserID == 0 {
		return errors.New("userID cannot be zero")
	}

	err := uc.blogRepo.Create(ctx, blog)

	if err != nil {
		return errors.New("failed to create blog")
	}
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

func (u *blogUsecase) DeleteBlog(ctx context.Context, ID int64, userID string) error {
	return u.blogRepo.DeleteByID(ctx, ID, userID)
}

func (uc *blogUsecase) GenerateBlogIdeas(topic string) (string, error) {
	return uc.aiService.GenerateBlogIdeas(topic)
}

func (uc *blogUsecase) SuggestBlogImprovements(content string) (string, error) {
	return uc.aiService.SuggestBlogImprovements(content)
}
func (uc *blogUsecase) GetAIService() domain.IAIService {
	return uc.aiService
}
func (uc *blogUsecase) FetchPaginatedBlogs(ctx context.Context, page, limit int) ([]*domain.Blog, int64, error) {
	return uc.blogRepo.FetchPaginatedBlogs(ctx, page, limit)

}

func (u *blogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog, tags []string) error { // NEW
	existing, err := u.blogRepo.FetchByID(ctx, blog.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("blog not found")
	}
	return u.blogRepo.Update(ctx, blog, tags)
}
