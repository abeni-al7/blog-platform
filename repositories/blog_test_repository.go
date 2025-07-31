package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abeni-al7/blog-platform/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)
	return gormDB, mock
}

func TestCreateBlog(t *testing.T) {
	db, mock := setupMockDB(t)

	repo := NewBlogRepository(db)

	ctx := context.Background()
	blog := &domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog content.",
		ID:      1,
	}
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "blogs"`).
		WithArgs(blog.Title, blog.Content, blog.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, blog)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
