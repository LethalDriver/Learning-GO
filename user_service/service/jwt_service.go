package service

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	expirationTimeString := os.Getenv("TOKEN_EXP_HS")
	if expirationTimeString == "" {
		return nil, errors.New("TOKEN_EXP_HS env variable not set")
	}

	// Parse the expiration time from string to integer
	expirationTimeHs, err := strconv.Atoi(expirationTimeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOKEN_EXP_HS: %v", err)
	}

	// Get the RSA private key
	privateKey, err := getRsaPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed getting rsa private key %v", err)
	}

	return &JwtService{
		expirationTimeHs: expirationTimeHs,
		privateKey:       privateKey,
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

// getRsaPrivateKey reads the RSA private key from the environment variable.
// It parses the RSA private key and returns the parsed private key.
func getRsaPrivateKey() (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile("private_key.pem")
	if err != nil {
		return nil, fmt.Errorf("error reading rsa private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}
	return privateKey, nil
}
