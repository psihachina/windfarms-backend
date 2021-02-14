package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

const salt = "nfnmskjfu73j4g5sdgv"
const singingKey = "dsaqwe123tjhcncj75#43GHD234"
const tokenTTL = 12 * time.Hour

// AuthService - ...
type AuthService struct {
	repo repository.Authorization
}

// tokenClaims - ...
type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

// NewAuthService - ...
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser - ...
func (s *AuthService) CreateUser(user models.User) error {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// GenerateToken - ...
func (s *AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.UserID,
	})

	return token.SignedString([]byte(singingKey))
}

// ParseToken -
func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid singing method")
		}

		return []byte(singingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

// generatePasswordHash - ...
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
