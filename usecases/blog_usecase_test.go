package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/abeni-al7/blog-platform/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBlogRepo struct {
	mock.Mock
}

func (m *mockBlogRepo) Create(ctx context.Context, blog *domain.Blog) error {
	args := m.Called(ctx, blog)
	return args.Error(0)
}
func TestBlogUsecase_CreateBlog_Success(t *testing.T) {
	mockRepo := new(mockBlogRepo)
	usecase := NewBlogUsecase(mockRepo)

	ctx := context.Background()
	blog := &domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog content.",
		ID:      2,
	}

	mockRepo.On("Create", ctx, blog).Return(nil)

	err := usecase.CreateBlog(ctx, blog)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBlogUsecase_CreateBlog_Error(t *testing.T) {
	mockRepo := new(mockBlogRepo)
	usecase := NewBlogUsecase(mockRepo)

	ctx := context.Background()
	blog := &domain.Blog{
		Title:   "Fail Blog",
		Content: "This blog will fail.",
		ID:      3,
	}

	mockRepo.On("Create", ctx, blog).Return(errors.New("failed to create blog"))

	err := usecase.CreateBlog(ctx, blog)
	assert.EqualError(t, err, "failed to create blog")
	mockRepo.AssertExpectations(t)

}
