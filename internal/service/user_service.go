package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	signUp       func(email, password string) (*domain.User, error)
	login        func(email, password string) (*domain.TokenPair, error)
	refreshToken func(refreshToken string) (*domain.TokenPair, error)
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	// Private helper for token generation, closed over by other functions if needed,
	// or kept as a private utility within the closure scope.
	// In this design, we can inject it or keep it internal.
	// Let's keep it as a shared private function for the factories.
	genToken := newTokenGenerator()

	return &userService{
		signUp:       newSignUpFunc(repo),
		login:        newLoginFunc(repo, genToken),
		refreshToken: newRefreshTokenFunc(repo, genToken),
	}
}

func (s *userService) SignUp(email, password string) (*domain.User, error) {
	return s.signUp(email, password)
}

func (s *userService) Login(email, password string) (*domain.TokenPair, error) {
	return s.login(email, password)
}

func (s *userService) RefreshToken(refreshToken string) (*domain.TokenPair, error) {
	return s.refreshToken(refreshToken)
}

// -------------------------------------------------------------------------
// Functional Implementations
// -------------------------------------------------------------------------

// TokenGeneratorFunc is a type for the token generation logic
type TokenGeneratorFunc func(userID uuid.UUID) (*domain.TokenPair, error)

func newTokenGenerator() TokenGeneratorFunc {
	return func(userID uuid.UUID) (*domain.TokenPair, error) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret"
		}

		// Access Token
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  userID.String(),
			"type": "access",
			"exp":  time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
		})
		accessTokenString, err := accessToken.SignedString([]byte(secret))
		if err != nil {
			return nil, err
		}

		// Refresh Token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  userID.String(),
			"type": "refresh",
			"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		})
		refreshTokenString, err := refreshToken.SignedString([]byte(secret))
		if err != nil {
			return nil, err
		}

		return &domain.TokenPair{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		}, nil
	}
}

func newSignUpFunc(repo domain.UserRepository) func(email, password string) (*domain.User, error) {
	return func(email, password string) (*domain.User, error) {
		// Check if user already exists
		if _, err := repo.FindByEmail(email); err == nil {
			return nil, errors.New("email already registered")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user := &domain.User{
			Email:    email,
			Password: string(hashedPassword),
		}

		if err := repo.Create(user); err != nil {
			return nil, err
		}

		return user, nil
	}
}

func newLoginFunc(repo domain.UserRepository, genToken TokenGeneratorFunc) func(email, password string) (*domain.TokenPair, error) {
	return func(email, password string) (*domain.TokenPair, error) {
		user, err := repo.FindByEmail(email)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, errors.New("invalid credentials")
		}

		return genToken(user.ID)
	}
}

func newRefreshTokenFunc(repo domain.UserRepository, genToken TokenGeneratorFunc) func(refreshToken string) (*domain.TokenPair, error) {
	return func(refreshTokenString string) (*domain.TokenPair, error) {
		token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "secret"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return nil, errors.New("invalid refresh token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		// Verify it's a refresh token
		if claims["type"] != "refresh" {
			return nil, errors.New("invalid token type")
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return nil, errors.New("invalid token subject")
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			return nil, errors.New("invalid user id in token")
		}

		return genToken(userID)
	}
}
