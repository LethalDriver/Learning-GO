package service

import (
    "context"
    "errors"
    "fmt"

    "example.com/chat_app/user_service/structs"
    "github.com/google/uuid"

    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
)

var (
    // ErrNoUser is returned when a user does not exist.
    ErrNoUser = errors.New("user doesn't exist")
    // ErrWrongPassword is returned when the provided password is incorrect.
    ErrWrongPassword = errors.New("incorrect password")
    // ErrUserExists is returned when a user already exists.
    ErrUserExists = errors.New("user already exists")
)

// UserRepository defines the interface for user data storage.
type UserRepository interface {
    // GetById retrieves a user by their ID.
    GetById(ctx context.Context, id string) (*structs.UserEntity, error)
    // GetByUsername retrieves a user by their username.
    GetByUsername(ctx context.Context, username string) (*structs.UserEntity, error)
    // Save stores a user entity in the repository.
    Save(ctx context.Context, user *structs.UserEntity) error
}

// UserService provides methods for user management.
type UserService struct {
    repo UserRepository
    jwt  *JwtService
}

// NewUserService creates a new UserService.
// It initializes the UserService with the provided UserRepository and JwtService.
func NewUserService(repo UserRepository, jwt *JwtService) *UserService {
    return &UserService{
        repo: repo,
        jwt:  jwt,
    }
}

// RegistrationRequest represents the data required to register a new user.
type RegistrationRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// LoginRequest represents the data required to log in a user.
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// GetUserDto retrieves a user DTO by their ID.
func (s *UserService) GetUserDto(ctx context.Context, userId string) (*structs.UserDto, error) {
    user, err := s.repo.GetById(ctx, userId)
    if err != nil {
        return nil, err
    }
    userDto := MapUserEntityToDto(user)
    return userDto, nil
}

// GetUser retrieves a user entity by their username.
func (s *UserService) GetUser(ctx context.Context, username string) (*structs.UserEntity, error) {
    user, err := s.repo.GetByUsername(ctx, username)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, ErrNoUser
        }
        return nil, fmt.Errorf("error getting user %q: %w", username, err)
    }
    return user, nil
}

// GetUserById retrieves a user entity by their ID.
func (s *UserService) GetUserById(ctx context.Context, id string) (*structs.UserEntity, error) {
    user, err := s.repo.GetById(ctx, id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, ErrNoUser
        }
        return nil, fmt.Errorf("error getting user %q: %w", id, err)
    }
    return user, nil
}

// RegisterUser registers a new user and returns the user DTO and JWT token.
// It validates the registration request, hashes the password, saves the user entity,
// and generates a JWT token.
func (s *UserService) RegisterUser(ctx context.Context, r RegistrationRequest) (*structs.UserDto, string, error) {
    err := s.validateRegistrationRequest(r)
    if err != nil {
        return nil, "", fmt.Errorf("registration request invalid: %w", err)
    }

    exists := s.checkIfUserExists(ctx, r.Username)
    if exists {
        return nil, "", ErrUserExists
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, "", fmt.Errorf("failed hashing password: %w", err)
    }

    user := &structs.UserEntity{
        Id:       uuid.New().String(),
        Username: r.Username,
        Email:    r.Email,
        Password: string(hashedPassword),
    }
    err = s.repo.Save(ctx, user)
    if err != nil {
        return nil, "", fmt.Errorf("failed saving user %q to the database: %w", user.Username, err)
    }

    token, err := s.jwt.GenerateToken(user.Id, user.Username)
    if err != nil {
        return nil, "", fmt.Errorf("failed generating jwt token: %w", err)
    }

    userDto := MapUserEntityToDto(user)

    return userDto, token, nil
}

// LoginUser logs in a user and returns the user DTO and JWT token.
// It retrieves the user entity by username, compares the hashed password,
// and generates a JWT token.
func (s *UserService) LoginUser(ctx context.Context, r LoginRequest) (*structs.UserDto, string, error) {
    user, err := s.repo.GetByUsername(ctx, r.Username)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, "", ErrNoUser
        }
        return nil, "", err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password))
    if err != nil {
        return nil, "", ErrWrongPassword
    }

    token, err := s.jwt.GenerateToken(user.Id, user.Username)
    if err != nil {
        return nil, "", fmt.Errorf("failed generating jwt token: %w", err)
    }

    userDto := MapUserEntityToDto(user)

    return userDto, token, nil
}

// validateRegistrationRequest validates the registration request data.
// It checks the email, username, and password for validity.
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

// checkIfUserExists checks if a user with the given username already exists.
func (s *UserService) checkIfUserExists(ctx context.Context, username string) bool {
    _, err := s.repo.GetByUsername(ctx, username)
    return err != mongo.ErrNoDocuments
}