package test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blog-platform/domain"
	"github.com/blog-platform/repositories"
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
	suite.repo = repositories.NewBlogRepository(gormDB)
}

func (suite *BlogRepoTestSuite) TestCreateBlog() {
	blog := &domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog content.",
		UserID:  1,
	}
	suite.mock.ExpectBegin()
	// GORM will insert all fields, so use AnyArg for those you don't care about
	suite.mock.ExpectQuery(`INSERT INTO "blogs"`).
		WithArgs(
			blog.Title,
			blog.Content,
			blog.ViewCount,
			blog.Likes,
			blog.Dislikes,
			blog.UserID,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	err := suite.repo.Create(context.Background(), blog)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *BlogRepoTestSuite) TestDeleteByID() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`DELETE FROM "blogs" WHERE id = \$1 AND user_id = \$2`).
		WithArgs(1, "user123").
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.DeleteByID(context.Background(), 1, "user123")
	assert.NoError(suite.T(), err)
}

func (suite *BlogRepoTestSuite) TestUpdateByID_Success() {
	updates := map[string]interface{}{
		"title":   "New Title",
		"content": "New Content",
	}
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "blogs" SET "content"=\$1,"title"=\$2,"updated_at"=\$3 WHERE id = \$4 AND user_id = \$5`).
		WithArgs(updates["content"], updates["title"], sqlmock.AnyArg(), 1, "user123").
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.UpdateByID(context.Background(), 1, "user123", updates)
	assert.NoError(suite.T(), err)
}

func (suite *BlogRepoTestSuite) TestUpdateByID_NotFound() {
	updates := map[string]interface{}{"title": "Nope"}
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "blogs" SET "title"=\$1,"updated_at"=\$2 WHERE id = \$3 AND user_id = \$4`).
		WithArgs(updates["title"], sqlmock.AnyArg(), 999, "user123").
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repo.UpdateByID(context.Background(), 999, "user123", updates)
	assert.Error(suite.T(), err)
}

func TestBlogRepoTestSuite(t *testing.T) {
	suite.Run(t, new(BlogRepoTestSuite))
}
