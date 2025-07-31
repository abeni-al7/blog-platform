package usecases

import (
	"context"

	"github.com/abeni-al7/blog-platform/domain"
)

type blogUsecase struct {
	blogRepo domain.IBlogRepository
}

func NewBlogUsecase(repo domain.IBlogRepository) domain.IBlogUsecase {
	return &blogUsecase{
		blogRepo: repo,
	}
}

func (uc blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog) error {
	return uc.blogRepo.Create(ctx, blog)
}
