package mocks

import "github.com/stretchr/testify/mock"

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(to []string, subject string, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}
