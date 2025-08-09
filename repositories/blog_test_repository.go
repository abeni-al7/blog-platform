package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blog-platform/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BlogRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo domain.IBlogRepository // Change from *BlogRepository to domain.IBlogRepository
}

func (suite *BlogRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(suite.T(), err)

	suite.db = gormDB
	suite.mock = mock
	suite.repo = NewBlogRepository(gormDB)
}

func (suite *BlogRepoTestSuite) TestCreateBlog() {
	blog := &domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog content.",
		ID:      1,
	}
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO "blogs"`).
		WithArgs(blog.Title, blog.Content, blog.ID, blog.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Create(context.Background(), blog)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *BlogRepoTestSuite) TestDeleteByID() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "blogs" SET "deleted_at"=`).
		WithArgs(sqlmock.AnyArg(), 1, "user123").
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.DeleteByID(context.Background(), 1, "user123")
	assert.NoError(suite.T(), err)
}

func TestBlogRepoTestSuite(t *testing.T) {
	suite.Run(t, new(BlogRepoTestSuite))
}
