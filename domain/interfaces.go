package domain

type IJWTInfrastructure interface {
	GenerateAccessToken(userID string, userRole string) (string, error)
	GenerateRefreshToken(userID string, userRole string) (string, error)
	ValidateAccessToken(authHeader string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}