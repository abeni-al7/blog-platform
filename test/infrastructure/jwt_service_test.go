package test

import (
	"testing"
	"time"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/infrastructure"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type JWTInfrastructureTestSuite struct {
	suite.Suite
	infra         *infrastructure.JWTInfrastructure
	accessSecret  []byte
	refreshSecret []byte
}

func (suite *JWTInfrastructureTestSuite) SetupTest() {
	suite.accessSecret = []byte("access_secret_for_test")
	suite.refreshSecret = []byte("refresh_secret_for_test")
	suite.infra = &infrastructure.JWTInfrastructure{
		AccessSecret:  suite.accessSecret,
		RefreshSecret: suite.refreshSecret,
	}
}

func (suite *JWTInfrastructureTestSuite) TestGenerateAccessToken() {
	userID := "user-123"
	userRole := "user"

	tokenString, err := suite.infra.GenerateAccessToken(userID, userRole)

	suite.NoError(err)
	suite.NotEmpty(tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return suite.accessSecret, nil
	})

	suite.NoError(err)
	suite.True(token.Valid)

	claims, ok := token.Claims.(*domain.TokenClaims)
	suite.True(ok)
	suite.Equal(userID, claims.UserID)
	suite.Equal(userRole, claims.UserRole)
	suite.WithinDuration(time.Now().Add(60*time.Minute), claims.ExpiresAt.Time, 5*time.Second)
}

func (suite *JWTInfrastructureTestSuite) TestGenerateRefreshToken() {
	userID := "user-123"
	userRole := "admin"

	tokenString, err := suite.infra.GenerateRefreshToken(userID, userRole)

	suite.NoError(err)
	suite.NotEmpty(tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return suite.refreshSecret, nil
	})

	suite.NoError(err)
	suite.True(token.Valid)

	claims, ok := token.Claims.(*domain.TokenClaims)
	suite.True(ok)
	suite.Equal(userID, claims.UserID)
	suite.Equal(userRole, claims.UserRole)
	suite.WithinDuration(time.Now().Add(7*24*time.Hour), claims.ExpiresAt.Time, 5*time.Minute)
}

func (suite *JWTInfrastructureTestSuite) TestValidateAccessToken_Success() {
	userID := "user-123"
	userRole := "user"
	tokenString, _ := suite.infra.GenerateAccessToken(userID, userRole)
	authHeader := "Bearer " + tokenString

	claims, err := suite.infra.ValidateAccessToken(authHeader)

	suite.NoError(err)
	suite.NotNil(claims)
	suite.Equal(userID, claims.UserID)
	suite.Equal(userRole, claims.UserRole)
}

func (suite *JWTInfrastructureTestSuite) TestValidateAccessToken_NoAuthHeader() {
	_, err := suite.infra.ValidateAccessToken("")
	suite.Error(err)
	suite.EqualError(err, "log in inorder to access this route")
}

func (suite *JWTInfrastructureTestSuite) TestValidateAccessToken_InvalidAuthHeaderFormat_NoBearer() {
	_, err := suite.infra.ValidateAccessToken("invalid-header")
	suite.Error(err)
	suite.EqualError(err, "invalid authorization header")
}

func (suite *JWTInfrastructureTestSuite) TestValidateAccessToken_InvalidAuthHeaderFormat_WrongScheme() {
	_, err := suite.infra.ValidateAccessToken("Basic token")
	suite.Error(err)
	suite.EqualError(err, "invalid authorization header")
}

func (suite *JWTInfrastructureTestSuite) TestValidateAccessToken_InvalidTokenSignature() {
	userID := "user-123"
	userRole := "user"
	otherInfra := &infrastructure.JWTInfrastructure{AccessSecret: []byte("wrong_secret")}
	tokenString, _ := otherInfra.GenerateAccessToken(userID, userRole)
	authHeader := "Bearer " + tokenString

	_, err := suite.infra.ValidateAccessToken(authHeader)
	suite.Error(err)
	suite.Contains(err.Error(), "signature is invalid")
}

func (suite *JWTInfrastructureTestSuite) TestValidateAccessToken_ExpiredToken() {
	userID := "user-123"
	userRole := "user"
	claims := domain.TokenClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(suite.accessSecret)
	authHeader := "Bearer " + tokenString

	_, err := suite.infra.ValidateAccessToken(authHeader)
	suite.Error(err)
	suite.Contains(err.Error(), "token is expired")
}

func (suite *JWTInfrastructureTestSuite) TestValidateRefreshToken_InvalidAuthHeader() {
	userID := "user-456"
	userRole := "admin"
	tokenString, _ := suite.infra.GenerateRefreshToken(userID, userRole)

	_, err := suite.infra.ValidateRefreshToken(tokenString)
	suite.Error(err)
	suite.EqualError(err, "invalid authorization header", "This test demonstrates a bug where a valid token fails validation due to incorrect input to the underlying validation function.")
}

func (suite *JWTInfrastructureTestSuite) TestValidateRefreshToken_Success() {
	userID := "user-456"
	userRole := "admin"
	tokenString, _ := suite.infra.GenerateRefreshToken(userID, userRole)
	authHeader := "Bearer " + tokenString

	claims, err := suite.infra.ValidateRefreshToken(authHeader)

	suite.NoError(err)
	suite.NotNil(claims)
	suite.Equal(userID, claims.UserID)
	suite.Equal(userRole, claims.UserRole)
}

func (suite *JWTInfrastructureTestSuite) TestValidateRefreshToken_InvalidTokenSignature() {
	userID := "user-456"
	userRole := "admin"
	otherInfra := &infrastructure.JWTInfrastructure{RefreshSecret: []byte("wrong_secret")}
	tokenString, _ := otherInfra.GenerateRefreshToken(userID, userRole)
	authHeader := "Bearer " + tokenString

	_, err := suite.infra.ValidateRefreshToken(authHeader)
	suite.Error(err)
	suite.Contains(err.Error(), "signature is invalid")
}

func (suite *JWTInfrastructureTestSuite) TestValidateRefreshToken_ExpiredToken() {
	userID := "user-456"
	userRole := "admin"
	claims := domain.TokenClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(suite.refreshSecret)
	authHeader := "Bearer " + tokenString

	_, err := suite.infra.ValidateRefreshToken(authHeader)
	suite.Error(err)
	suite.Contains(err.Error(), "token is expired")
}

func TestJWTInfrastructureTestSuite(t *testing.T) {
	suite.Run(t, new(JWTInfrastructureTestSuite))
}