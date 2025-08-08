package usecases

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strconv"
	"unicode"

	"github.com/blog-platform/domain"
)

type UserUsecase struct {
	userRepo        domain.IUserRepository
	emailService    domain.IEmailInfrastructure
	passwordService domain.IPasswordInfrastructure
	jwtService      domain.IJWTInfrastructure
	tokenRepo       domain.ITokenRepository
}

func NewUserUsecase(ur domain.IUserRepository, es domain.IEmailInfrastructure, ps domain.IPasswordInfrastructure, js domain.IJWTInfrastructure, tr domain.ITokenRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:        ur,
		emailService:    es,
		passwordService: ps,
		jwtService:      js,
		tokenRepo:       tr,
	}
}

func (uu *UserUsecase) Register(user *domain.User) (domain.User, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return domain.User{}, errors.New("missing required fields")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return domain.User{}, errors.New("invalid email format")
	}

	if !uu.validatePassword(user.Password) {
		return domain.User{}, errors.New("password must be consisted of at least one uppercase character, one lowercase character, one punctuation character, one number and be at least of length 8")
	}

	_, err = uu.userRepo.FetchByUsername(user.Username)
	if err == nil {
		return domain.User{}, errors.New("this username is already in use")
	}

	_, err = uu.userRepo.FetchByEmail(user.Email)
	if err == nil {
		return domain.User{}, errors.New("this email is already in use")
	}

	user.Status = "inactive"
	user.Password, err = uu.passwordService.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, errors.New(err.Error())
	}

	registeredUser, err := uu.userRepo.Register(user)
	if err != nil {
		return domain.User{}, errors.New("unable to register user")
	}

	emailContent := fmt.Sprintf("%v://%v:%v/user/%v/activate", os.Getenv("PROTOCOL"), os.Getenv("DOMAIN"), os.Getenv("PORT"), registeredUser.ID)
	err = uu.emailService.SendEmail([]string{registeredUser.Email}, "Activate Account", emailContent)
	if err != nil {
		return domain.User{}, errors.New("unable to send activation link")
	}

	return registeredUser, nil
}

func (uu *UserUsecase) Login(identifier string, password string) (string, string, error) {
	user, err := uu.userRepo.FetchByUsername(identifier)
	if err != nil {
		_, err := mail.ParseAddress(identifier)
		if err != nil {
			return "", "", errors.New("invalid email format")
		}

		user, err = uu.userRepo.FetchByEmail(identifier)
		if err != nil {
			return "", "", errors.New("invalid identifier")
		}
	}

	if !uu.validatePassword(password) {
		return "", "", errors.New("invalid password format")
	}
	err = uu.passwordService.ComparePassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := uu.jwtService.GenerateAccessToken(strconv.FormatInt(user.ID, 10), user.Role)
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	refreshToken, err := uu.jwtService.GenerateRefreshToken(strconv.FormatInt(user.ID, 10), user.Role)
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	accessTokenObj := domain.Token{
		Type:    "access",
		Content: accessToken,
		Status:  "active",
		UserID:  user.ID,
	}
	refreshTokenObj := domain.Token{
		Type:    "refresh",
		Content: refreshToken,
		Status:  "active",
		UserID:  user.ID,
	}

	err = uu.tokenRepo.Save(&accessTokenObj)
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	err = uu.tokenRepo.Save(&refreshTokenObj)
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	return accessToken, refreshToken, nil
}

func (uu *UserUsecase) RefreshToken(authHeader string) (string, string, error) {
	claims, err := uu.jwtService.ValidateRefreshToken(authHeader)
	if err != nil {
		return "", "", err
	}

	accessToken, err := uu.jwtService.GenerateAccessToken(claims.UserID, claims.UserRole)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uu.jwtService.GenerateRefreshToken(claims.UserID, claims.UserRole)
	if err != nil {
		return "", "", err
	}

	// persist new tokens
	uid, _ := strconv.ParseInt(claims.UserID, 10, 64)
	accessTokenObj := domain.Token{Type: "access", Content: accessToken, Status: "active", UserID: uid}
	refreshTokenObj := domain.Token{Type: "refresh", Content: refreshToken, Status: "active", UserID: uid}
	if err = uu.tokenRepo.Save(&accessTokenObj); err != nil {
		return "", "", err
	}
	if err = uu.tokenRepo.Save(&refreshTokenObj); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uu *UserUsecase) validatePassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func (uu *UserUsecase) ActivateAccount(id string) error {
	_, err := uu.userRepo.Fetch(id)
	if err != nil {
		return err
	}

	err = uu.userRepo.ActivateAccount(id)
	if err != nil {
		return err
	}

	return nil
}

func (uu UserUsecase) GetUserProfile(userID int64) (*domain.User, error) {
	user, err := uu.userRepo.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu UserUsecase) UpdateUserProfile(userID int64, updates map[string]interface{}) error {
	return uu.userRepo.UpdateUserProfile(userID, updates)
}

func (uu *UserUsecase) ResetPassword(userID string, oldPassword string, newPassword string) error {
	user, err := uu.userRepo.Fetch(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := uu.passwordService.ComparePassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	if !uu.validatePassword(newPassword) {
		return errors.New("password must be consisted of at least one uppercase character, one lowercase character, one punctuation character, one number and be at least of length 8")
	}

	hashed, err := uu.passwordService.HashPassword(newPassword)
	if err != nil {
		return errors.New("could not hash password")
	}

	if err := uu.userRepo.ResetPassword(userID, hashed); err != nil {
		return errors.New("could not update password")
	}
	return nil
}
