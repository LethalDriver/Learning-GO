package service

import (
	"errors"
	"fmt"

	"example.com/myproject/entity"
	"example.com/myproject/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoUser = errors.New("user doesn't exist")
	ErrWrongPassword = errors.New("incorrect password")
	ErrUserExists = errors.New("user already exists")
)

type UserService struct {
	repo      repository.UserRepository
	jwt *JwtService
}

func NewUserService(repo repository.UserRepository, jwt *JwtService) *UserService {
	return &UserService{
		repo: repo,
		jwt: jwt,
	}
}

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserService) RegisterUser(r RegistrationRequest) error {
	err := s.validateRegistrationRequest(r)
	if err != nil {
		return err
	}
	exists := s.checkIfUserExists(r.Username)
	if exists {
		return ErrUserExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := entity.NewUserEntity(r.Username, r.Email, string(hashedPassword))
	err = s.repo.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) LoginUser(r LoginRequest) (string, error) {
	user, err := s.repo.GetByUsername(r.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", ErrNoUser
		}
		return "", errors.New("unknown error while logging in")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password))
	if err != nil {
		return "", ErrWrongPassword
	}

	token, err := s.jwt.GenerateToken(user.Id, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed generating jwt token: %w", err)
	}

	return token, nil 
}

func (s *UserService) validateRegistrationRequest(r RegistrationRequest) error {
	err := ValidateEmail(r.Email)
	if err != nil {
		return err
	}
	err = ValidateUsername(r.Username)
	if err != nil {
		return err
	}
	err = ValidatePassword(r.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) checkIfUserExists(username string) bool {
	_, err := s.repo.GetByUsername(username)
	return err != mongo.ErrNoDocuments
}
