package test

import (
	"context"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/test/mocks"
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
	mockRepo *mocks.MockBlogRepo
	mockAI   *MockAIService
	usecase  domain.IBlogUsecase
}

func (suite *BlogUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockBlogRepo)
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

func (suite *BlogUsecaseTestSuite) TestAddComment_Success() {
	ctx := context.Background()
	expected := &domain.Comment{ID: 1, BlogID: 10, UserID: 5, Content: "Nice!"}
	suite.mockRepo.On("AddComment", ctx, int64(10), int64(5), "Nice!").Return(expected, nil)
	c, err := suite.usecase.AddComment(ctx, 10, 5, "Nice!")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, c)
}

func (suite *BlogUsecaseTestSuite) TestAddComment_Validation() {
	ctx := context.Background()
	_, err := suite.usecase.AddComment(ctx, 0, 5, "Hi")
	assert.Error(suite.T(), err)
	_, err = suite.usecase.AddComment(ctx, 10, 0, "Hi")
	assert.Error(suite.T(), err)
	_, err = suite.usecase.AddComment(ctx, 10, 5, "  ")
	assert.Error(suite.T(), err)
}

func (suite *BlogUsecaseTestSuite) TestGetComments_Success() {
	ctx := context.Background()
	list := []*domain.Comment{{ID: 1}, {ID: 2}}
	suite.mockRepo.On("ListComments", ctx, int64(10), 1, 10).Return(list, int64(2), nil)
	got, total, err := suite.usecase.GetComments(ctx, 10, 1, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Equal(suite.T(), list, got)
}

func (suite *BlogUsecaseTestSuite) TestGetComments_InvalidBlog() {
	ctx := context.Background()
	_, _, err := suite.usecase.GetComments(ctx, 0, 1, 10)
	assert.Error(suite.T(), err)
}

func TestBlogUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BlogUsecaseTestSuite))
}
