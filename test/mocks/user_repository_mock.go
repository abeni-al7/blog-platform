package mocks

import (
	"github.com/blog-platform/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Register(user *domain.User) (domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) FetchByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) FetchByEmail(email string) (domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) ActivateAccount(idStr string) error {
	args := m.Called(idStr)
	return args.Error(0)
}

func (m *MockUserRepository) Fetch(idStr string) (domain.User, error) {
	args := m.Called(idStr)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserProfile(userID int64) (*domain.User, error) {
	args := m.Called(userID)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) UpdateUserProfile(userID int64, updates map[string]interface{}) error {
	args := m.Called(userID, updates)
	return args.Error(0)
}
func (m *MockUserRepository) ResetPassword(idStr string, newPassword string) error {
	args := m.Called(idStr, newPassword)
	return args.Error(0)
}
