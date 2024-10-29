package middleware

import (
	"net/http"
	"strings"

	"github.com/rhajizada/donezo/internal/auth"
)

// AuthMiddleware returns a middleware function that validates JWT tokens.
func AuthMiddleware(secret []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, auth.BearerPrefix)
			if tokenString == authHeader {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Validate the token
			err := auth.ValidateToken(secret, tokenString)
			if err != nil {
				// Differentiate between expired tokens and other validation errors
				if err.Error() == "token has expired" {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				}
				if err.Error() == "token not valid yet" {
					http.Error(w, "Token not valid yet", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Token is valid; proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
