package mocks

import (
	"github.com/blog-platform/domain"
	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
    mock.Mock
}

func (m *MockTokenRepository) FetchByContent(content string) (domain.Token, error) {
    args := m.Called(content)
    return args.Get(0).(domain.Token), args.Error(1)
}

func (m *MockTokenRepository) Save(token *domain.Token) error {
	args := m.Called(token)
	return args.Error(0)
}