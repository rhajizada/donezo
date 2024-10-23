package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(secret []byte, issuer string, expiration time.Duration) (string, error) {
	// Define the standard claims
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
	}

	// Create the token with the HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
