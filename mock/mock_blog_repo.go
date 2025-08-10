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

func (m *MockBlogRepo) FetchByID(ctx context.Context, id int64) (*domain.Blog, error) {
	args := m.Called(ctx, id)
	if blog, ok := args.Get(0).(*domain.Blog); ok {
		return blog, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepo) FetchAll(ctx context.Context) ([]*domain.Blog, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Blog), args.Error(1)
}

func (m *MockBlogRepo) FetchPaginatedBlogs(ctx context.Context, page, limit int) ([]*domain.Blog, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*domain.Blog), args.Get(1).(int64), args.Error(2)
}

func (m *MockBlogRepo) IncrementView(ctx context.Context, blogID int64) error {
	args := m.Called(ctx, blogID)
	return args.Error(0)
}

func (m *MockBlogRepo) AddLike(ctx context.Context, blogID int64, userID int64) error {
	args := m.Called(ctx, blogID, userID)
	return args.Error(0)
}

func (m *MockBlogRepo) RemoveLike(ctx context.Context, blogID int64, userID int64) error {
	args := m.Called(ctx, blogID, userID)
	return args.Error(0)
}

func (m *MockBlogRepo) GetPopularity(ctx context.Context, blogID int64) (int, int, error) {
	args := m.Called(ctx, blogID)
	return args.Int(0), args.Int(1), args.Error(2)
}
