package test

import (
	"errors"
	"net/smtp"
	"testing"

	"github.com/blog-platform/infrastructure"
	"github.com/stretchr/testify/suite"
)

type EmailServiceTestSuite struct {
	suite.Suite
	emailService *infrastructure.SMTPEmailService
}

func (suite *EmailServiceTestSuite) SetupTest() {
	suite.emailService = &infrastructure.SMTPEmailService{
		Host:     "smtp.example.com",
		Port:     "587",
		Username: "user",
		Password: "password",
		From:     "from@example.com",
	}
}

func (suite *EmailServiceTestSuite) TestSendEmail_Success() {
	suite.emailService.SendMailFn = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return nil
	}

	err := suite.emailService.SendEmail([]string{"to@example.com"}, "Test Subject", "Test Body")
	suite.NoError(err)
}

func (suite *EmailServiceTestSuite) TestSendEmail_Failure() {
	suite.emailService.SendMailFn = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return errors.New("failed to send email")
	}

	err := suite.emailService.SendEmail([]string{"to@example.com"}, "Test Subject", "Test Body")
	suite.Error(err)
	suite.Equal("failed to send email", err.Error())
}

func TestEmailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmailServiceTestSuite))
}

