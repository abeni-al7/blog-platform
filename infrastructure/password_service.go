package infrastructure

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordInfrastructure struct {}

func NewPasswordInfrastructure() *PasswordInfrastructure {
	return &PasswordInfrastructure{}
}

func (infra *PasswordInfrastructure) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("unable to hash password")
	}
	return string(hashedPassword), nil
}

func (infra *PasswordInfrastructure) ComparePassword(correctPassword []byte, inputPassword []byte) error {
	if bcrypt.CompareHashAndPassword(correctPassword, inputPassword) != nil {
		return errors.New("invalid credentials")
	}
	return nil
}