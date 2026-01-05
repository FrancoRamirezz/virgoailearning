package utils

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/config"
	"backend/models"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenInvalid = errors.New("token is invalid")

// GetJWTSecret returns the JWT secret from environment variables
func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

// GetUserIDFromToken extracts user ID from JWT token in request context
func GetUserIDFromToken(r *http.Request) (uint, error) {
	user := r.Context().Value("user")
	if user == nil {
		return 0, ErrTokenInvalid
	}

	// Cast to the actual user model type
	userStruct, ok := user.(models.User)
	if !ok {
		return 0, ErrTokenInvalid
	}

	return userStruct.ID, nil
}

// GetUserIDFromClaims extracts user ID from JWT claims
func GetUserIDFromClaims(claims jwt.Claims) (uint, error) {
	if claimsMap, ok := claims.(jwt.MapClaims); ok {
		if userIDFloat, ok := claimsMap["user_id"].(float64); ok {
			return uint(userIDFloat), nil
		}
		if userIDStr, ok := claimsMap["user_id"].(string); ok {
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				return 0, err
			}
			return uint(userID), nil
		}
	}
	return 0, ErrTokenInvalid
}

// GenerateToken creates a JWT token for a user
func GenerateToken(userID uint, email, role string, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(cfg.JWT.Expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}
