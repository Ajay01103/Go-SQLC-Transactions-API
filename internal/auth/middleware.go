package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	jsonutil "github.com/Ajay01103/goTransactonsAPI/internal/json"
)

// contextKey is a private type for context keys in the auth package.
type contextKey string

// ContextKeyUserID is the key used to store the authenticated userID in the request context.
const ContextKeyUserID contextKey = "userID"

// RequireAuth returns a middleware that validates the Bearer JWT and injects
// the userID into the request context. Rejects requests with 401 if the token
// is missing, malformed, or expired.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				jsonutil.Write(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				jsonutil.Write(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
