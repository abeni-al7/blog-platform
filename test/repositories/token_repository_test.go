package test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blog-platform/domain"
	"github.com/blog-platform/repositories"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TokenRepositoryTestSuite struct {
    suite.Suite
    db   *sql.DB
    mock sqlmock.Sqlmock
    repo *repositories.TokenRepository
}

func (s *TokenRepositoryTestSuite) SetupTest() {
    var err error
    s.db, s.mock, err = sqlmock.New()
    s.Require().NoError(err)

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: s.db,
    }), &gorm.Config{})
    s.Require().NoError(err)

    s.repo = repositories.NewTokenRepository(gormDB)
}

func (s *TokenRepositoryTestSuite) TearDownTest() {
	s.mock.ExpectationsWereMet()
}

func (s *TokenRepositoryTestSuite) TestFetchByContent_Success() {
    expectedQuery := `SELECT * FROM "tokens" WHERE content = $1 AND "tokens"."deleted_at" IS NULL ORDER BY "tokens"."id" LIMIT $2`

    expectedToken := domain.Token{
        ID:      1,
        Content: "test_token",
        UserID:  3,
        Status:  "active",
    }

    rows := sqlmock.NewRows([]string{"id", "content", "user_id", "status"}).
        AddRow(expectedToken.ID, expectedToken.Content, expectedToken.UserID, expectedToken.Status)

    s.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
        WithArgs("test_token", 1).
        WillReturnRows(rows)

    token, err := s.repo.FetchByContent("test_token")

    s.NoError(err)
    s.Equal(expectedToken, token)
}

func (s *TokenRepositoryTestSuite) TestFetchByContent_NotFound() {
    expectedQuery := `SELECT * FROM "tokens" WHERE content = $1 AND "tokens"."deleted_at" IS NULL ORDER BY "tokens"."id" LIMIT $2`

    s.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
        WithArgs("non_existent_token", 1).
        WillReturnError(gorm.ErrRecordNotFound)

    _, err := s.repo.FetchByContent("non_existent_token")

    s.Error(err)
}

func (s *TokenRepositoryTestSuite) TestFetchByContent_DBError() {
    expectedQuery := `SELECT * FROM "tokens" WHERE content = $1 AND "tokens"."deleted_at" IS NULL ORDER BY "tokens"."id" LIMIT $2`
    dbError := errors.New("some db error")

    s.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
        WithArgs("any_token", 1).
        WillReturnError(dbError)

    _, err := s.repo.FetchByContent("any_token")

    s.Error(err)
}

func (s *TokenRepositoryTestSuite) TestSave_Success() {
	tokenToSave := &domain.Token{
        Type: "access",
		Content: "new_token",
		UserID:  1,
		Status:  "active",
	}

	expectedQuery := `INSERT INTO "tokens" ("created_at","updated_at","deleted_at","type","content","status","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","id"`

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tokenToSave.Type, tokenToSave.Content, tokenToSave.Status, tokenToSave.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := s.repo.Save(tokenToSave)

	s.NoError(err)
}

func (s *TokenRepositoryTestSuite) TestSave_DBError() {
    tokenToSave := &domain.Token{
        Content: "new_token",
        UserID:  1,
        Status:  "active",
    }
    dbError := errors.New("some db error")

    expectedQuery := `INSERT INTO "tokens" ("created_at","updated_at","deleted_at","type","content","status","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","id"`

    s.mock.ExpectBegin()
    s.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
        WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tokenToSave.Type, tokenToSave.Content, tokenToSave.Status, tokenToSave.UserID).
        WillReturnError(dbError)
    s.mock.ExpectRollback()

    err := s.repo.Save(tokenToSave)

    s.Error(err)
}

// helper to expect soft delete (GORM sets only deleted_at in this model)
func (s *TokenRepositoryTestSuite) expectSoftDelete(content string, retErr error) {
    expectedExec := `UPDATE "tokens" SET "deleted_at"=$1 WHERE content = $2 AND "tokens"."deleted_at" IS NULL`
    s.mock.ExpectBegin()
    exec := s.mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
        WithArgs(sqlmock.AnyArg(), content)
    if retErr != nil {
        exec.WillReturnError(retErr)
        s.mock.ExpectRollback()
    } else {
        exec.WillReturnResult(sqlmock.NewResult(0, 1))
        s.mock.ExpectCommit()
    }
}

func (s *TokenRepositoryTestSuite) TestDelete_Success() {
    s.expectSoftDelete("del_token", nil)
    err := s.repo.Delete("del_token")
    s.NoError(err)
}

func (s *TokenRepositoryTestSuite) TestDelete_DBError() {
    dbErr := errors.New("delete failed")
    s.expectSoftDelete("bad_token", dbErr)
    err := s.repo.Delete("bad_token")
    s.Error(err)
    s.Equal("delete failed", err.Error())
}

func TestTokenRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(TokenRepositoryTestSuite))
}