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
func (suite *BlogUsecaseTestSuite) TestFetchBlogsByFilter_Success() {
	ctx := context.Background()
	filter := domain.BlogFilter{TitleContains: "Test"}
	expectedBlogs := []*domain.Blog{
		{ID: 1, Title: "Test Blog", Content: "Content"},
	}

	suite.mockRepo.On("FetchByFilter", ctx, filter).Return(expectedBlogs, nil)

	blogs, err := suite.usecase.FetchBlogsByFilter(ctx, filter)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedBlogs, blogs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *BlogUsecaseTestSuite) TestFetchBlogsByFilter_Error() {
	ctx := context.Background()
	filter := domain.BlogFilter{TitleContains: "Fail"}

	suite.mockRepo.On("FetchByFilter", ctx, filter).Return([]*domain.Blog(nil), assert.AnError)

	blogs, err := suite.usecase.FetchBlogsByFilter(ctx, filter)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), blogs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestBlogUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BlogUsecaseTestSuite))
}
