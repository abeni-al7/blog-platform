package infrastructure

import (
	"errors"
	"strings"
	"time"

	"github.com/blog-platform/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTInfrastructure struct {
	AccessSecret []byte
	RefreshSecret []byte
}

func NewJWTInfrastructure(accessSecret, refreshSecret []byte) (*JWTInfrastructure, error) {
	if len(accessSecret) == 0 || len(refreshSecret) == 0 {
		return nil, errors.New("access and refresh secrets cannot be empty")
	}
	return &JWTInfrastructure{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
	}, nil
}

func (infra *JWTInfrastructure) GenerateAccessToken(userID string, userRole string) (string, error) {
	if userID == "" || userRole == "" {
		return "", errors.New("userID and userRole cannot be empty")
	}

	claims := domain.TokenClaims{
		UserID:    userID,
		UserRole:  userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(infra.AccessSecret)
}

func (infra *JWTInfrastructure) GenerateRefreshToken(userID string, userRole string) (string, error) {
	if userID == "" || userRole == "" {
		return "", errors.New("userID and userRole cannot be empty")
	}

	claims := domain.TokenClaims{
		UserID: userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(infra.RefreshSecret)
}

func (infra *JWTInfrastructure) validateToken(authHeader string, secret []byte) (*domain.TokenClaims, error) {
	if authHeader == "" {
		return &domain.TokenClaims{}, errors.New("log in inorder to access this route")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		return &domain.TokenClaims{}, errors.New("invalid authorization header")
	}

	tokenString := authParts[1]

	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	exp, err := claims.GetExpirationTime() 
	if err != nil {
		return nil, errors.New("token validation failed")
	}

	if exp.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func (infra *JWTInfrastructure) ValidateAccessToken(authHeader string) (*domain.TokenClaims, error) {
	return infra.validateToken(authHeader, infra.AccessSecret)
}

func (infra *JWTInfrastructure) ValidateRefreshToken(authHeader string) (*domain.TokenClaims, error) {
	return infra.validateToken(authHeader, infra.RefreshSecret)
}