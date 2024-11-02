package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type ContextKey string

const (
	UserIdKey   ContextKey = "userId"
	UsernameKey ContextKey = "username"
)

type AuthService struct {
	publicKey *rsa.PublicKey
}

// JWTMiddleware is the authorization middleware
func JWTMiddleware(authService *AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Validate the token
		claims, err := authService.ValidateTokenWithClaims(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-Id", userId)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// ValidateTokenWithClaims validates the token and returns the claims if the token is valid.
func (s *AuthService) ValidateTokenWithClaims(tokenString string) (jwt.MapClaims, error) {
	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Ensure the token's signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("failed to extract claims")
}

// getRsaPublicKey reads the RSA public key from the environment variable.
// It parses the RSA public key and returns the parsed public key.
func getRsaPublicKey() (*rsa.PublicKey, error) {
	pemData, err := os.ReadFile("public_key.pem")
	if err != nil {
		return nil, fmt.Errorf("error reading rsa private key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %v", err)
	}
	return publicKey, nil
}
