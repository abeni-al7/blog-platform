package test

import (
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

type UserRepositoryTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	repo *repositories.UserRepository
}

func (s *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	s.Require().NoError(err)

	s.mock = mock
	s.repo = repositories.NewUserRepository(s.DB)
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	s.mock.ExpectationsWereMet()
}

func (s *UserRepositoryTestSuite) TestRegister_Success() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		Status:   "inactive",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","password","role","bio","profile_picture","phone","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, user.Email, user.Password, "", "", "", "", user.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	createdUser, err := s.repo.Register(user)
	s.NoError(err)
	s.NotZero(createdUser.ID)
	s.Equal(user.Username, createdUser.Username)
}

func (s *UserRepositoryTestSuite) TestRegister_Error() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		Status:   "inactive",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","password","role","bio","profile_picture","phone","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, user.Email, user.Password, "", "", "", "", user.Status).
		WillReturnError(errors.New("db error"))
	s.mock.ExpectRollback()

	_, err := s.repo.Register(user)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestFetchByEmail_Success() {
	email := "test@example.com"
	user := domain.User{ID: 1, Email: email, Username: "testuser"}

	rows := sqlmock.NewRows([]string{"id", "email", "username"}).
		AddRow(user.ID, user.Email, user.Username)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	foundUser, err := s.repo.FetchByEmail(email)
	s.NoError(err)
	s.Equal(user, foundUser)
}

func (s *UserRepositoryTestSuite) TestFetchByEmail_NotFound() {
	email := "test@example.com"

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := s.repo.FetchByEmail(email)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestFetchByUsername_Success() {
	username := "testuser"
	user := domain.User{ID: 1, Email: "test@example.com", Username: username}

	rows := sqlmock.NewRows([]string{"id", "email", "username"}).
		AddRow(user.ID, user.Email, user.Username)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(username, 1).
		WillReturnRows(rows)

	foundUser, err := s.repo.FetchByUsername(username)
	s.NoError(err)
	s.Equal(user, foundUser)
}

func (s *UserRepositoryTestSuite) TestFetchByUsername_NotFound() {
	username := "testuser"

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(username).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := s.repo.FetchByUsername(username)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestActivateAccount_Success() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "users"."deleted_at" IS NULL`)).
		WithArgs("active", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.ActivateAccount("1")
	s.NoError(err)
}

func (s *UserRepositoryTestSuite) TestActivateAccount_InvalidID() {
	err := s.repo.ActivateAccount("invalid")
	s.Error(err)
	s.EqualError(err, "invalid id")
}

func (s *UserRepositoryTestSuite) TestActivateAccount_DBError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "status"=$1 WHERE "id" = $2`)).
		WithArgs("active", 1).
		WillReturnError(errors.New("db error"))
	s.mock.ExpectRollback()

	err := s.repo.ActivateAccount("1")
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestActivateAccount_NoRowsAffected() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "status"=$1 WHERE "id" = $2`)).
		WithArgs("active", 1).
		WillReturnResult(sqlmock.NewResult(1, 0))
	s.mock.ExpectCommit()

	err := s.repo.ActivateAccount("1")
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestFetch_Success() {
	user := domain.User{ID: 1, Email: "test@example.com", Username: "testuser"}

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "role", "bio", "profile_picture", "phone", "status"}).
		AddRow(user.ID, user.Email, user.Username, user.Password, user.Role, user.Bio, user.ProfilePicture, user.Phone, user.Status)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	foundUser, err := s.repo.Fetch("1")
	s.NoError(err)
	s.Equal(user, foundUser)
}

func (s *UserRepositoryTestSuite) TestFetch_InvalidID() {
	_, err := s.repo.Fetch("invalid")
	s.Error(err)
	s.EqualError(err, "invalid id")
}

func (s *UserRepositoryTestSuite) TestFetch_NotFound() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := s.repo.Fetch("1")
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestGetUserProfile_Success() {
	user := domain.User{ID: 1, Email: "test@example.com", Username: "testuser"}
	rows := sqlmock.NewRows([]string{"id", "email", "username"}).
		AddRow(user.ID, user.Email, user.Username)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(user.ID, 1).
		WillReturnRows(rows)

	foundUser, err := s.repo.GetUserProfile(user.ID)
	s.NoError(err)
	s.Equal(user.ID, foundUser.ID)
	s.Equal(user.Email, foundUser.Email)
	s.Equal(user.Username, foundUser.Username)
}

func (s *UserRepositoryTestSuite) TestGetUserProfile_NotFound() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(int64(2), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := s.repo.GetUserProfile(2)
	s.NoError(err)
	s.Nil(user)
}

func (s *UserRepositoryTestSuite) TestUpdateUserProfile_Success() {
	userID := int64(1)
	updates := map[string]interface{}{
		"Username": "updateduser",
		"Bio":      "updated bio",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"bio"=$2 WHERE id = $3 AND "users"."deleted_at" IS NULL`)).
		WithArgs("updateduser", "updated bio", userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateUserProfile(userID, updates)
	s.NoError(err)
}

func (s *UserRepositoryTestSuite) TestUpdateUserProfile_NoFields() {
	userID := int64(1)
	updates := map[string]interface{}{
		"NotAllowed": "value",
	}
	err := s.repo.UpdateUserProfile(userID, updates)
	s.NoError(err)
}
func (s *UserRepositoryTestSuite) TestResetPassword_Success() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(1, "old_hashed"))

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password"=$1,"updated_at"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs("new_hashed", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.repo.ResetPassword("1", "new_hashed")
	s.NoError(err)
}

func (s *UserRepositoryTestSuite) TestResetPassword_InvalidID() {
	err := s.repo.ResetPassword("abc", "new_hashed")
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestResetPassword_UserNotFound() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	err := s.repo.ResetPassword("1", "new_hashed")
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestResetPassword_UpdateError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(1, "old_hashed"))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password"=$1,"updated_at"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs("new_hashed", sqlmock.AnyArg(), 1).
		WillReturnError(errors.New("db error"))
	s.mock.ExpectRollback()
	err := s.repo.ResetPassword("1", "new_hashed")
	s.Error(err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
