package usecases

import (
	"context"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BlogUsecaseTestSuite struct {
	suite.Suite
	mockRepo *mock.MockBlogRepo
	usecase  domain.IBlogUsecase
}

func (suite *BlogUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(mock.MockBlogRepo)
	suite.usecase = NewBlogUsecase(suite.mockRepo)
}

func (suite *BlogUsecaseTestSuite) TestCreateBlog_Success() {
	ctx := context.Background()
	blog := &domain.Blog{
		ID:      1,
		Title:   "Test Blog",
		Content: "This is a test blog content.",
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
	}

	// Simulate repo error
	suite.mockRepo.On("Create", ctx, blog).Return(assert.AnError)

	tags := []string{}
	err := suite.usecase.CreateBlog(ctx, blog, tags)
	assert.EqualError(suite.T(), err, "failed to create blog")
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestBlogUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BlogUsecaseTestSuite))
}
