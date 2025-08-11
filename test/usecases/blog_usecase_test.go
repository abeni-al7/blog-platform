package test

import (
	"context"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/mock"
	"github.com/blog-platform/usecases"
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
	suite.usecase = usecases.NewBlogUsecase(suite.mockRepo, suite.mockAI)
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
	updates := map[string]interface{}{"Title": "New", "Content": "Body"}
	suite.mockRepo.On("UpdateByID", ctx, int64(1), "123", updates).Return(nil)
	err := suite.usecase.UpdateBlog(ctx, 1, "123", updates)
	assert.NoError(suite.T(), err)
}

func (suite *BlogUsecaseTestSuite) TestUpdateBlog_InvalidID() {
	ctx := context.Background()
	err := suite.usecase.UpdateBlog(ctx, 0, "123", map[string]interface{}{"Title": "X"})
	assert.Error(suite.T(), err)
}

func (suite *BlogUsecaseTestSuite) TestUpdateBlog_EmptyTitle() {
	ctx := context.Background()
	err := suite.usecase.UpdateBlog(ctx, 1, "123", map[string]interface{}{"Title": ""})
	assert.Error(suite.T(), err)
}

func TestBlogUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BlogUsecaseTestSuite))
}
