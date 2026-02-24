package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 7 * 24 * time.Hour // 7 days

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// generateToken creates a signed HS256 JWT for the given user with a 7-day TTL.
func generateToken(userID, email, secret string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signed, nil
}
