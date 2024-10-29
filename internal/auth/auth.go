// package auth
//
// import (
// 	"time"
//
// 	"github.com/golang-jwt/jwt/v4"
// )
//
// func GenerateJWT(secret []byte, issuer string, expiration time.Duration) (string, error) {
// 	// Define the standard claims
// 	claims := jwt.RegisteredClaims{
// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
// 	}
//
// 	// Create the token with the HS256 signing method and the claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
// 	// Sign the token with the secret
// 	tokenString, err := token.SignedString(secret)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return tokenString, nil
// }
//
// func ValidateJWT(token string) error {
// 	return nil
// }

package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	Issuer       = "donezo"
	BearerPrefix = "Bearer "
)

// Claims defines the structure of JWT claims
type Claims struct {
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token with the given secret and expiration duration.
func GenerateToken(secret []byte, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates the JWT token string using the provided secret.
// It returns an error if the token is invalid, expired, or has an incorrect issuer.
func ValidateToken(secret []byte, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and specifically HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		// Check if the error is due to token expiration
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return errors.New("token has expired")
			}
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return errors.New("token not valid yet")
			}
		}
		return err // Return the original error for other cases
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Verify the issuer
		if claims.Issuer != Issuer {
			return errors.New("invalid issuer")
		}

		// Explicitly check token expiration
		if claims.ExpiresAt == nil {
			return errors.New("token does not have an expiration time")
		}

		if time.Until(claims.ExpiresAt.Time) <= 0 {
			return errors.New("token has expired")
		}

		return nil // Token is valid and not expired
	}

	return errors.New("invalid token")
}
