package user

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"tasklybe/pkg/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrValidation         = errors.New("validation error")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrJWTSecretMissing   = errors.New("jwt secret missing")
)

type authClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateAPIKey(u *User) (string, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return "", ErrJWTSecretMissing
	}

	now := time.Now()
	claims := authClaims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func RegisterUser(input RegisterUserRequest) (*User, error) {

	var existing User
	err := db.DB.First(&existing, "email = ?", input.Email).Error
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := User{
		ID:       uuid.NewString(),
		Email:    input.Email,
		Password: string(hashed),
		Name:     strings.TrimSpace(input.Name),
	}

	if err := db.DB.Create(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func LoginUser(input LoginUserRequest) (*User, string, error) {

	var u User
	if err := db.DB.First(&u, "email = ?", input.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	apiKey, err := generateAPIKey(&u)
	if err != nil {
		return nil, "", err
	}

	return &u, apiKey, nil
}

func GetUserByID(id string) (*User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", ErrValidation)
	}

	var u User
	if err := db.DB.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func ToUserResponse(u *User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
		Name:      u.Name,
	}
}
