package usecases

import (
	"context"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MockAIService struct{}

func (m *MockAIService) GenerateBlogIdeas(topic string) (string, error) {
	return "Mocked blog ideas", nil
}
func (m *MockAIService) SuggestBlogImprovements(content string) (string, error) {
	return "Mocked improvement", nil
}

type BlogUsecaseTestSuite struct {
	suite.Suite
	mockRepo *mock.MockBlogRepo
	mockAI   *MockAIService
	usecase  domain.IBlogUsecase
}

func (suite *BlogUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(mock.MockBlogRepo)
	suite.mockAI = &MockAIService{}
	suite.usecase = NewBlogUsecase(suite.mockRepo, suite.mockAI)
}

func (suite *BlogUsecaseTestSuite) TestCreateBlog_Success() {
	ctx := context.Background()
	blog := &domain.Blog{
		ID:      1,
		Title:   "Test Blog",
		Content: "This is a test blog content.",
		UserID:  123,
	}

	suite.mockRepo.On("Create", ctx, blog).Return(nil)

	tags := []string{}
	err := suite.usecase.CreateBlog(ctx, blog, tags)
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *BlogUsecaseTestSuite) TestCreateBlogError() {
	ctx := context.Background()
	blog := &domain.Blog{
		ID:      2,
		Title:   "Fail Blog",
		Content: "This blog will fail.",
		UserID:  456,
	}

	// Simulate repo error
	suite.mockRepo.On("Create", ctx, blog).Return(assert.AnError)

	tags := []string{}
	err := suite.usecase.CreateBlog(ctx, blog, tags)
	assert.EqualError(suite.T(), err, "failed to create blog")
	suite.mockRepo.AssertExpectations(suite.T())
}
func (suite *BlogUsecaseTestSuite) TestUpdateBlog_Success() {
	ctx := context.Background()
	blog := &domain.Blog{
		ID:      1,
		Title:   "Updated Title",
		Content: "Updated content",
		UserID:  123,
	}

	tags := []string{"go", "programming"}
	// Expect FetchByID call to verify blog exists (optional but good practice)
	suite.mockRepo.On("FetchByID", ctx, blog.ID).Return(blog, nil)
	suite.mockRepo.On("Update", ctx, blog, tags).Return(nil)

	// For tags, mock FindOrCreateTag and LinkTagToBlog calls accordingly
	suite.mockRepo.On("FindOrCreateTag", ctx, "go").Return(int64(1), nil)
	suite.mockRepo.On("LinkTagToBlog", ctx, blog.ID, int64(1)).Return(nil)
	suite.mockRepo.On("FindOrCreateTag", ctx, "programming").Return(int64(2), nil)
	suite.mockRepo.On("LinkTagToBlog", ctx, blog.ID, int64(2)).Return(nil)

	err := suite.usecase.UpdateBlog(ctx, blog, tags)
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *BlogUsecaseTestSuite) TestUpdateBlog_NotFound() {
	ctx := context.Background()
	blog := &domain.Blog{
		ID:      99,
		Title:   "Nonexistent",
		Content: "Content",
		UserID:  123,
	}

	suite.mockRepo.On("FetchByID", ctx, blog.ID).Return(nil, assert.AnError)
	err := suite.usecase.UpdateBlog(ctx, blog, []string{})
	assert.Error(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestBlogUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BlogUsecaseTestSuite))
}
