package mocks

import (
	"github.com/blog-platform/domain"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID string, userRole string) (string, error) {
	args := m.Called(userID, userRole)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID string, userRole string) (string, error) {
	args := m.Called(userID, userRole)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateAccessToken(authHeader string) (*domain.TokenClaims, error) {
	args := m.Called(authHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(authHeader string) (*domain.TokenClaims, error) {
	args := m.Called(authHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}
