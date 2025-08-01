package mock

import (
	"context"

	"github.com/blog-platform/domain"
	"github.com/stretchr/testify/mock"
)

type MockBlogRepo struct {
	mock.Mock
}

func (m *MockBlogRepo) Create(ctx context.Context, blog *domain.Blog) error {
	args := m.Called(ctx, blog)
	return args.Error(0)
}

func (m *MockBlogRepo) FindOrCreateTag(ctx context.Context, tag string) (int64, error) {
	args := m.Called(ctx, tag)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBlogRepo) LinkTagToBlog(ctx context.Context, blogID int64, tagID int64) error {
	args := m.Called(ctx, blogID, tagID)
	return args.Error(0)
}
