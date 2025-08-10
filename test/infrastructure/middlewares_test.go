package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	mockJWTService *mocks.MockJWTService
	middleware     *infrastructure.Middleware
	router         *gin.Engine
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.mockJWTService = new(mocks.MockJWTService)
	suite.middleware = infrastructure.NewMiddleware(suite.mockJWTService)
	suite.router = gin.Default()
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_Success() {
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	claims := &domain.TokenClaims{UserID: "user123", UserRole: "user"}
	suite.mockJWTService.On("ValidateAccessToken", "Bearer valid_token").Return(claims, nil)

	suite.router.GET("/test", suite.middleware.AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockJWTService.AssertExpectations(suite.T())
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_NoAuthHeader() {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	suite.mockJWTService.On("ValidateAccessToken", "").Return(nil, errors.New("no auth header"))

	suite.router.GET("/test", suite.middleware.AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"invalid token"}`, w.Body.String())
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	suite.mockJWTService.On("ValidateAccessToken", "Bearer invalid_token").Return(nil, errors.New("invalid token"))

	suite.router.GET("/test", suite.middleware.AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"invalid token"}`, w.Body.String())
	suite.mockJWTService.AssertExpectations(suite.T())
}

func (suite *MiddlewareTestSuite) TestAdminMiddleware_Success() {
	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	suite.router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	}, suite.middleware.AdminMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *MiddlewareTestSuite) TestAdminMiddleware_Forbidden() {
	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	suite.router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	}, suite.middleware.AdminMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.JSONEq(suite.T(), `{"error":"unauthorized to access this route"}`, w.Body.String())
}

func (suite *MiddlewareTestSuite) TestAdminMiddleware_NoRole() {
	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	suite.router.GET("/admin", suite.middleware.AdminMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.JSONEq(suite.T(), `{"error":"unauthorized to access this route"}`, w.Body.String())
}

func (suite *MiddlewareTestSuite) TestAccountOwnerMiddleware_Success() {
	req, _ := http.NewRequest("GET", "/users/user123", nil)
	w := httptest.NewRecorder()

	suite.router.GET("/users/:id", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, suite.middleware.AccountOwnerMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *MiddlewareTestSuite) TestAccountOwnerMiddleware_Forbidden() {
	req, _ := http.NewRequest("GET", "/users/user456", nil)
	w := httptest.NewRecorder()

	suite.router.GET("/users/:id", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, suite.middleware.AccountOwnerMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.JSONEq(suite.T(), `{"error":"unauthorized to access this route"}`, w.Body.String())
}

func (suite *MiddlewareTestSuite) TestAccountOwnerMiddleware_NoUserID() {
	req, _ := http.NewRequest("GET", "/users/user123", nil)
	w := httptest.NewRecorder()

	suite.router.GET("/users/:id", suite.middleware.AccountOwnerMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.JSONEq(suite.T(), `{"error":"unauthorized to access this route"}`, w.Body.String())
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}