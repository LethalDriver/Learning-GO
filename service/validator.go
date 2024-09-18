package service

import (
	"errors"
	"regexp"
	"unicode"
)

// Validator is an interface that defines methods for validating email, username, and password.
type Validator interface {
	ValidatePassword(password string) error
	ValidateUsername(username string) error
	ValidateEmail(email string) error
}

// ValidatorImpl is an implementation of the Validator interface.
type ValidatorImpl struct {
}

// ValidateEmail validates the given email address using a regular expression pattern.
// It returns an error if the email format is invalid.
func (v *ValidatorImpl) ValidateEmail(email string) error {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegexPattern)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateUsername validates the length of the given username.
// It returns an error if the username length is less than 3 or greater than 30 characters.
func (v *ValidatorImpl) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}
	return nil
}

// ValidatePassword validates the length and composition of the given password.
// It returns an error if the password length is less than 8 characters or if it does not contain at least one uppercase letter, one lowercase letter, one digit, and one special character.
func (v *ValidatorImpl) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}

	return nil
}