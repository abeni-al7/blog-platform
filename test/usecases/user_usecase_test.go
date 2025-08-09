package test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/test/mocks"
	"github.com/blog-platform/usecases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo     *mocks.MockUserRepository
	emailService *mocks.MockEmailService
	pwdService   *mocks.MockPasswordService
	jwtService   *mocks.MockJWTService
	tokenRepo    *mocks.MockTokenRepository
	userUsecase  domain.IUserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepo = new(mocks.MockUserRepository)
	suite.emailService = new(mocks.MockEmailService)
	suite.pwdService = new(mocks.MockPasswordService)
	suite.jwtService = new(mocks.MockJWTService)
	suite.tokenRepo = new(mocks.MockTokenRepository)
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, suite.jwtService, suite.tokenRepo)
	os.Setenv("PROTOCOL", "http")
	os.Setenv("DOMAIN", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("EMAIL_SENDER", "test@example.com")
}

func (suite *UserUsecaseTestSuite) TestRegister_Success() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	createdUser := *user
	createdUser.ID = 1

	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.pwdService.On("HashPassword", user.Password).Return("hashedpassword", nil)
	suite.userRepo.On("Register", mock.AnythingOfType("*domain.User")).Return(createdUser, nil)
	suite.emailService.On("SendEmail", []string{user.Email}, "Activate Account", "http://localhost:8080/user/1/activate").Return(nil)

	_, err := suite.userUsecase.Register(user)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestRegister_MissingFields() {
	user := &domain.User{}
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("missing required fields", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_InvalidEmail() {
	user := &domain.User{
		Username: "testuser",
		Email:    "invalid-email",
		Password: "Password123!",
	}
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("invalid email format", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_InvalidPassword() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "weak",
	}
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("password must be consisted of at least one uppercase character, one lowercase character, one punctuation character, one number and be at least of length 8", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_UsernameInUse() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, nil)
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("this username is already in use", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_EmailInUse() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, nil)
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("this email is already in use", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_RegistrationFails_HashError() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.pwdService.On("HashPassword", user.Password).Return("", errors.New("hash error"))

	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestRegister_RegistrationFails_RepoError() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.pwdService.On("HashPassword", user.Password).Return("hashedpassword", nil)
	suite.userRepo.On("Register", mock.AnythingOfType("*domain.User")).Return(domain.User{}, errors.New("db error"))

	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("unable to register user", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_SendEmailFails() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	createdUser := *user
	createdUser.ID = 1

	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.pwdService.On("HashPassword", user.Password).Return("hashedpassword", nil)
	suite.userRepo.On("Register", mock.AnythingOfType("*domain.User")).Return(createdUser, nil)
	suite.emailService.On("SendEmail", []string{user.Email}, "Activate Account", "http://localhost:8080/user/1/activate").Return(errors.New("email error"))

	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("unable to send activation link", err.Error())
}

func (suite *UserUsecaseTestSuite) TestActivateAccount_Success() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, nil)
	suite.userRepo.On("ActivateAccount", "1").Return(nil)
	err := suite.userUsecase.ActivateAccount("1")
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestActivateAccount_UserNotFound() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, errors.New("not found"))
	err := suite.userUsecase.ActivateAccount("1")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestActivateAccount_ActivationFails() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, nil)
	suite.userRepo.On("ActivateAccount", "1").Return(errors.New("db error"))
	err := suite.userUsecase.ActivateAccount("1")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_Success() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("refresh_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil).Twice()

	accessToken, refreshToken, err := suite.userUsecase.Login("testuser", "Password123!")

	suite.NoError(err)
	suite.Equal("access_token", accessToken)
	suite.Equal("refresh_token", refreshToken)
	suite.tokenRepo.AssertNumberOfCalls(suite.T(), "Save", 2)
}

func (suite *UserUsecaseTestSuite) TestLogin_InvalidIdentifier() {
	suite.userRepo.On("FetchByUsername", "unknown").Return(domain.User{}, errors.New("not found"))
	_, _, err := suite.userUsecase.Login("unknown", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_InvalidPassword() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("WrongPassword!")).Return(errors.New("wrong password"))
	_, _, err := suite.userUsecase.Login("testuser", "WrongPassword!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_GenerateAccessTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("", errors.New("jwt error"))

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_GenerateRefreshTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("", errors.New("jwt error"))

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_SaveAccessTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("refresh_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db error")).Once()

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_SaveRefreshTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("refresh_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil).Once()
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db error")).Once()

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestGetUserProfile_Success() {
	expectedUser := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	suite.userRepo.On("GetUserProfile", int64(1)).Return(expectedUser, nil)
	user, err := suite.userUsecase.GetUserProfile(1)
	suite.NoError(err)
	suite.Equal(expectedUser, user)
}

func (suite *UserUsecaseTestSuite) TestGetUserProfile_NotFound() {
	suite.userRepo.On("GetUserProfile", int64(2)).Return(nil, nil)
	user, err := suite.userUsecase.GetUserProfile(2)
	suite.NoError(err)
	suite.Nil(user)
}

func (suite *UserUsecaseTestSuite) TestUpdateUserProfile_Success() {
	userID := int64(1)
	updates := map[string]interface{}{
		"Username": "updateduser",
		"Bio":      "updated bio",
	}
	suite.userRepo.On("UpdateUserProfile", userID, updates).Return(nil)
	err := suite.userUsecase.UpdateUserProfile(userID, updates)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestUpdateUserProfile_NoFields() {
	userID := int64(1)
	updates := map[string]interface{}{}
	suite.userRepo.On("UpdateUserProfile", userID, updates).Return(nil)
	err := suite.userUsecase.UpdateUserProfile(userID, updates)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestRefreshToken_Success() {
	jwtMock := new(mocks.MockJWTService)
	tokenMock := new(mocks.MockTokenRepository)
	suite.jwtService = jwtMock
	suite.tokenRepo = tokenMock
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, jwtMock, tokenMock)
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	authHeader := "Bearer old_refresh"
	jwtMock.On("ValidateRefreshToken", authHeader).Return(claims, nil)
	jwtMock.On("GenerateAccessToken", "1", "user").Return("new_access", nil)
	jwtMock.On("GenerateRefreshToken", "1", "user").Return("new_refresh", nil)
	tokenMock.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil).Twice()
	access, refresh, err := suite.userUsecase.RefreshToken(authHeader)
	suite.NoError(err)
	suite.Equal("new_access", access)
	suite.Equal("new_refresh", refresh)
}

func (suite *UserUsecaseTestSuite) TestRefreshToken_ValidateError() {
	jwtMock := new(mocks.MockJWTService)
	tokenMock := new(mocks.MockTokenRepository)
	suite.jwtService = jwtMock
	suite.tokenRepo = tokenMock
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, jwtMock, tokenMock)
	authHeader := "Bearer bad"
	jwtMock.On("ValidateRefreshToken", authHeader).Return((*domain.TokenClaims)(nil), errors.New("invalid token"))
	_, _, err := suite.userUsecase.RefreshToken(authHeader)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestRefreshToken_GenerateAccessTokenError() {
	jwtMock := new(mocks.MockJWTService)
	tokenMock := new(mocks.MockTokenRepository)
	suite.jwtService = jwtMock
	suite.tokenRepo = tokenMock
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, jwtMock, tokenMock)
	authHeader := "Bearer old_refresh"
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	jwtMock.On("ValidateRefreshToken", authHeader).Return(claims, nil)
	jwtMock.On("GenerateAccessToken", "1", "user").Return("", errors.New("gen err"))
	_, _, err := suite.userUsecase.RefreshToken(authHeader)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestRefreshToken_GenerateRefreshTokenError() {
	jwtMock := new(mocks.MockJWTService)
	tokenMock := new(mocks.MockTokenRepository)
	suite.jwtService = jwtMock
	suite.tokenRepo = tokenMock
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, jwtMock, tokenMock)
	authHeader := "Bearer old_refresh"
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	jwtMock.On("ValidateRefreshToken", authHeader).Return(claims, nil)
	jwtMock.On("GenerateAccessToken", "1", "user").Return("new_access", nil)
	jwtMock.On("GenerateRefreshToken", "1", "user").Return("", errors.New("gen err"))
	_, _, err := suite.userUsecase.RefreshToken(authHeader)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestRefreshToken_SaveAccessTokenError() {
	jwtMock := new(mocks.MockJWTService)
	tokenMock := new(mocks.MockTokenRepository)
	suite.jwtService = jwtMock
	suite.tokenRepo = tokenMock
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, jwtMock, tokenMock)
	authHeader := "Bearer old_refresh"
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	jwtMock.On("ValidateRefreshToken", authHeader).Return(claims, nil)
	jwtMock.On("GenerateAccessToken", "1", "user").Return("new_access", nil)
	jwtMock.On("GenerateRefreshToken", "1", "user").Return("new_refresh", nil)
	tokenMock.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db err")).Once()
	_, _, err := suite.userUsecase.RefreshToken(authHeader)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestRefreshToken_SaveRefreshTokenError() {
	jwtMock := new(mocks.MockJWTService)
	tokenMock := new(mocks.MockTokenRepository)
	suite.jwtService = jwtMock
	suite.tokenRepo = tokenMock
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, jwtMock, tokenMock)
	authHeader := "Bearer old_refresh"
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	jwtMock.On("ValidateRefreshToken", authHeader).Return(claims, nil)
	jwtMock.On("GenerateAccessToken", "1", "user").Return("new_access", nil)
	jwtMock.On("GenerateRefreshToken", "1", "user").Return("new_refresh", nil)
	tokenMock.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil).Once()
	tokenMock.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db err")).Once()
	_, _, err := suite.userUsecase.RefreshToken(authHeader)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestResetPassword_Success() {
	user := domain.User{ID: 1, Password: "old_hashed"}
	suite.userRepo.On("Fetch", "1").Return(user, nil)
	suite.pwdService.On("ComparePassword", []byte("old_hashed"), []byte("OldPass123!")).Return(nil)
	suite.pwdService.On("HashPassword", "NewPass123!").Return("new_hashed", nil)
	suite.userRepo.On("ResetPassword", "1", "new_hashed").Return(nil)
	err := suite.userUsecase.ResetPassword("1", "OldPass123!", "NewPass123!")
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestResetPassword_UserNotFound() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, errors.New("not found"))
	err := suite.userUsecase.ResetPassword("1", "OldPass123!", "NewPass123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestResetPassword_InvalidOldPassword() {
	user := domain.User{ID: 1, Password: "old_hashed"}
	suite.userRepo.On("Fetch", "1").Return(user, nil)
	suite.pwdService.On("ComparePassword", []byte("old_hashed"), []byte("WrongOld123!")).Return(errors.New("mismatch"))
	err := suite.userUsecase.ResetPassword("1", "WrongOld123!", "NewPass123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestResetPassword_InvalidNewPasswordFormat() {
	user := domain.User{ID: 1, Password: "old_hashed"}
	suite.userRepo.On("Fetch", "1").Return(user, nil)
	suite.pwdService.On("ComparePassword", []byte("old_hashed"), []byte("OldPass123!")).Return(nil)
	err := suite.userUsecase.ResetPassword("1", "OldPass123!", "weak")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestResetPassword_HashError() {
	user := domain.User{ID: 1, Password: "old_hashed"}
	suite.userRepo.On("Fetch", "1").Return(user, nil)
	suite.pwdService.On("ComparePassword", []byte("old_hashed"), []byte("OldPass123!")).Return(nil)
	suite.pwdService.On("HashPassword", "NewPass123!").Return("", errors.New("hash fail"))
	err := suite.userUsecase.ResetPassword("1", "OldPass123!", "NewPass123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestResetPassword_UpdateError() {
	user := domain.User{ID: 1, Password: "old_hashed"}
	suite.userRepo.On("Fetch", "1").Return(user, nil)
	suite.pwdService.On("ComparePassword", []byte("old_hashed"), []byte("OldPass123!")).Return(nil)
	suite.pwdService.On("HashPassword", "NewPass123!").Return("new_hashed", nil)
	suite.userRepo.On("ResetPassword", "1", "new_hashed").Return(errors.New("db error"))
	err := suite.userUsecase.ResetPassword("1", "OldPass123!", "NewPass123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_Success() {
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	suite.jwtService.On("ValidateAccessToken", "Bearer token123").Return(claims, nil)
	suite.pwdService.On("HashPassword", "NewPass123!").Return("new_hashed", nil)
	suite.userRepo.On("ResetPassword", "1", "new_hashed").Return(nil)
	err := suite.userUsecase.UpdatePasswordDirect("1", "NewPass123!", "token123")
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_MissingToken() {
	err := suite.userUsecase.UpdatePasswordDirect("1", "NewPass123!", "")
	suite.Error(err)
	suite.Equal("token required", err.Error())
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_InvalidToken() {
	suite.jwtService.On("ValidateAccessToken", "Bearer badtoken").Return((*domain.TokenClaims)(nil), errors.New("invalid token"))
	err := suite.userUsecase.UpdatePasswordDirect("1", "NewPass123!", "badtoken")
	suite.Error(err)
	suite.Equal("invalid or expired token", err.Error())
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_TokenUserMismatch() {
	claims := &domain.TokenClaims{UserID: "2", UserRole: "user"}
	suite.jwtService.On("ValidateAccessToken", "Bearer tokenMismatch").Return(claims, nil)
	err := suite.userUsecase.UpdatePasswordDirect("1", "NewPass123!", "tokenMismatch")
	suite.Error(err)
	suite.Equal("token does not match user", err.Error())
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_InvalidPasswordFormat() {
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	suite.jwtService.On("ValidateAccessToken", "Bearer tokenFormat").Return(claims, nil)
	err := suite.userUsecase.UpdatePasswordDirect("1", "weak", "tokenFormat")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_HashError() {
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	suite.jwtService.On("ValidateAccessToken", "Bearer tokenHash").Return(claims, nil)
	suite.pwdService.On("HashPassword", "NewPass123!").Return("", errors.New("hash fail"))
	err := suite.userUsecase.UpdatePasswordDirect("1", "NewPass123!", "tokenHash")
	suite.Error(err)
	suite.Equal("could not hash password", err.Error())
}

func (suite *UserUsecaseTestSuite) TestUpdatePasswordDirect_UpdateError() {
	claims := &domain.TokenClaims{UserID: "1", UserRole: "user"}
	suite.jwtService.On("ValidateAccessToken", "Bearer tokenUpdate").Return(claims, nil)
	suite.pwdService.On("HashPassword", "NewPass123!").Return("new_hashed", nil)
	suite.userRepo.On("ResetPassword", "1", "new_hashed").Return(errors.New("db error"))
	err := suite.userUsecase.UpdatePasswordDirect("1", "NewPass123!", "tokenUpdate")
	suite.Error(err)
	suite.Equal("could not update password", err.Error())
}

func (suite *UserUsecaseTestSuite) TestForgotPassword_Success() {
	user := domain.User{ID: 1, Email: "user@example.com", Role: "user"}
	suite.userRepo.On("FetchByEmail", user.Email).Return(user, nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("reset_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil)
	suite.emailService.On("SendEmail", []string{user.Email}, "Reset Password", mock.MatchedBy(func(body string) bool {
		return strings.Contains(body, "/password/1/update?token=reset_token")
	})).Return(nil)

	err := suite.userUsecase.ForgotPassword(user.Email)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestForgotPassword_EmptyEmail() {
	err := suite.userUsecase.ForgotPassword("")
	suite.Error(err)
	suite.Equal("email required", err.Error())
}

func (suite *UserUsecaseTestSuite) TestForgotPassword_UserNotFound() {
	suite.userRepo.On("FetchByEmail", "missing@example.com").Return(domain.User{}, errors.New("not found"))
	err := suite.userUsecase.ForgotPassword("missing@example.com")
	suite.Error(err)
	suite.Equal("user not found", err.Error())
}

func (suite *UserUsecaseTestSuite) TestForgotPassword_GenerateTokenError() {
	user := domain.User{ID: 1, Email: "user@example.com", Role: "user"}
	suite.userRepo.On("FetchByEmail", user.Email).Return(user, nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("", errors.New("gen err"))
	err := suite.userUsecase.ForgotPassword(user.Email)
	suite.Error(err)
	suite.Equal("could not generate reset token", err.Error())
	suite.tokenRepo.AssertNotCalled(suite.T(), "Save", mock.Anything)
	suite.emailService.AssertNotCalled(suite.T(), "SendEmail", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *UserUsecaseTestSuite) TestForgotPassword_PersistTokenError() {
	user := domain.User{ID: 1, Email: "user@example.com", Role: "user"}
	suite.userRepo.On("FetchByEmail", user.Email).Return(user, nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("reset_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db err"))
	err := suite.userUsecase.ForgotPassword(user.Email)
	suite.Error(err)
	suite.Equal("could not persist reset token", err.Error())
	suite.emailService.AssertNotCalled(suite.T(), "SendEmail", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *UserUsecaseTestSuite) TestForgotPassword_SendEmailError() {
	user := domain.User{ID: 1, Email: "user@example.com", Role: "user"}
	suite.userRepo.On("FetchByEmail", user.Email).Return(user, nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("reset_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil)
	suite.emailService.On("SendEmail", []string{user.Email}, "Reset Password", mock.AnythingOfType("string")).Return(errors.New("smtp err"))
	err := suite.userUsecase.ForgotPassword(user.Email)
	suite.Error(err)
	suite.Equal("could not send reset link", err.Error())
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
