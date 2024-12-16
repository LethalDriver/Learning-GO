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

// AuthService provides methods for JWT token validation.
type AuthService struct {
	publicKey *rsa.PublicKey
}

// JWTMiddleware is the authorization middleware
// It reads the JWT token from the Authorization header, validates the token, and extracts the user ID from the token claims.
// The user ID is then appended to X-User-Id header in the request.
func JWTMiddleware(authService *AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

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

		next.ServeHTTP(w, r)
	})
}

// JWTQueryMiddleware reads the JWT token from the query parameter, validates it, and appends the user ID to the X-User-Id header.
func JWTQueryMiddleware(authService *AuthService, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.URL.Query().Get("token")
        if tokenString == "" {
            http.Error(w, "Token query parameter missing", http.StatusUnauthorized)
            return
        }

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

        next.ServeHTTP(w, r)
    })
}

// ValidateTokenWithClaims validates the token and returns the claims if the token is valid.
func (s *AuthService) ValidateTokenWithClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

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
