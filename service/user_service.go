package service

import (
	"errors"

	"example.com/myproject/entity"
	"example.com/myproject/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      repository.UserRepository
	validator Validator
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type RegistrationRequest struct {
	Username string
	Email    string
	Password string
}

func (s *UserService) RegisterUser(r RegistrationRequest) error {
	err := s.validateRegistrationRequest(r)
	if err != nil {
		return err
	}
	exists := s.checkIfUserExists(r.Username)
	if exists {
		return errors.New("user already exists")
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

func (s *UserService) validateRegistrationRequest(r RegistrationRequest) error {
	err := s.validator.ValidateEmail(r.Email)
	if err != nil {
		return err
	}
	err = s.validator.ValidateUsername(r.Username)
	if err != nil {
		return err
	}
	err = s.validator.ValidatePassword(r.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) checkIfUserExists(username string) bool {
	_, err := s.repo.GetByUsername(username)
	return err != mongo.ErrNoDocuments
}
