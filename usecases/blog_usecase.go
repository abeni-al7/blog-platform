package usecases

import (
	"context"
	"errors"

	"github.com/blog-platform/domain"
	"gorm.io/gorm"
)

type blogUsecase struct {
	blogRepo domain.IBlogRepository
}

func NewBlogUsecase(repo domain.IBlogRepository) domain.IBlogUsecase {
	return &blogUsecase{
		blogRepo: repo,
	}
}

func (uc blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, tags []string, userID int64) error {
	// prevent empty strings from being added
	blog.UserID = userID
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

	for _, tag := range tags {
		if tag == "" {
			continue // skip empty tags
		}
		tagID, err := uc.blogRepo.FindOrCreateTag(ctx, tag)
		if err != nil {
			return err
		}
		err = uc.blogRepo.LinkTagToBlog(ctx, int64(blog.ID), tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *blogUsecase) DeleteBlog(ctx context.Context, ID int64, userID int64) error {
	blog, err := uc.blogRepo.FetchByID(ctx, ID)
	// Validate userID
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog not found")
		}
		return err
	}

	if blog.UserID != userID {
		return errors.New("unauthorized to delete this blog")
	}
	if userID == 0 {
		return errors.New("userID cannot be zero")
	}

	// Call the repository method to delete the blog
	err = uc.blogRepo.DeleteBlog(ctx, ID)
	if err != nil {
		return errors.New("failed to delete blog")
	}

	return uc.blogRepo.DeleteBlog(ctx, ID)
}
