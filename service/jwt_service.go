package service

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtService represents a JWT service.
type JwtService struct {
	expirationTimeHs int
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
}

// NewJwtService creates a new instance of JwtService.
// It initializes the JwtService with the expiration time, private key, and public key.
// The expiration time is read from the TOKEN_EXPIRATION_HS environment variable.
// The private key is read from the RSA_PRIVATE_KEY environment variable.
// The public key is read from the RSA_PUBLIC_KEY environment variable.
func NewJwtService() (*JwtService, error) {
	// Read the expiration time from the environment variable
	expirationTimeString := os.Getenv("TOKEN_EXPIRATION_HS")
	if expirationTimeString == "" {
		return nil, errors.New("TOKEN_EXPIRATION_HS env variable not set")
	}

	// Parse the expiration time from string to integer
	expirationTimeHs, err := strconv.Atoi(expirationTimeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOKEN_EXPIRATION_HS: %v", err)
	}

	// Get the RSA private key
	privateKey, err := getRsaPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed getting rsa private key from env variables: %v", err)
	}

	// Get the RSA public key
	publicKey, err := getRsaPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed getting rsa public key from env variables: %v", err)
	}

	return &JwtService{
		expirationTimeHs: expirationTimeHs,
		privateKey:       privateKey,
		publicKey:        publicKey,
	}, nil
}

// GenerateToken generates a JWT token for the given user ID and username.
// It uses the RSA private key to sign the token.
// The expiration time of the token is set based on the expiration time in hours.
func (s *JwtService) GenerateToken(userId string, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expirationTimeHs))),
	})

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func (s *JwtService) ValidateToken(tokenString string) error {
    // Load the RSA public key
    publicKey, err := getRsaPublicKey()
    if err != nil {
        return fmt.Errorf("failed to get RSA public key: %v", err)
    }

    // Parse and verify the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Ensure the token's signing method is RSA
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return publicKey, nil
    })

    if err != nil {
        return fmt.Errorf("failed to parse token: %v", err)
    }

    // Check if the token is valid
    if !token.Valid {
        return errors.New("invalid token")
    }

    return nil
}


// getRsaPublicKey reads the RSA public key from the environment variable.
// It parses the RSA public key and returns the parsed public key.
func getRsaPublicKey() (*rsa.PublicKey, error) {
	publicKeyPEM := os.Getenv("RSA_PUBLIC_KEY")
	if publicKeyPEM == "" {
		return nil, errors.New("RSA_PUBLIC_KEY env variable not set")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %v", err)
	}
	return publicKey, nil
}

// getRsaPrivateKey reads the RSA private key from the environment variable.
// It parses the RSA private key and returns the parsed private key.
func getRsaPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyPEM := os.Getenv("RSA_PRIVATE_KEY")
	if privateKeyPEM == "" {
		return nil, errors.New("RSA_PRIVATE_KEY env variable not set")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}
	return privateKey, nil
}

// AuthMiddleware is the authorization middleware
func AuthMiddleware(jwtService *JwtService, next http.Handler) http.Handler {
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
        err := jwtService.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        // Call the next handler
        next.ServeHTTP(w, r)
    })
}


