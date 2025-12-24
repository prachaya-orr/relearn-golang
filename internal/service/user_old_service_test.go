package service_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"github.com/prachaya-orr/relearn-golang/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// Reuse MockUserRepository from user_service_test.go if possible,
// but since it's likely in the same package `service_test`, it should be available?
// Wait, `user_service_test.go` package is `service_test`.
// So `MockUserRepository` is already defined in the package scope.
// I can just reuse it!

func TestUserOldService_SignUp(t *testing.T) {
	repo := NewMockUserRepo()
	svc := service.NewUserOldService(repo)

	t.Run("Success", func(t *testing.T) {
		email := "test_old@example.com"
		password := "password123"

		user, err := svc.SignUp(email, password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Email != email {
			t.Errorf("expected email %s, got %s", email, user.Email)
		}
		if user.ID == uuid.Nil {
			t.Error("expected ID to be set")
		}
		if user.Password == password {
			t.Error("expected password to be hashed")
		}
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		email := "dup_old@example.com"
		svc.SignUp(email, "pass")

		_, err := svc.SignUp(email, "pass2")
		if err == nil {
			t.Error("expected error for duplicate email")
		}
		if err.Error() != "email already registered" {
			t.Errorf("expected 'email already registered', got '%v'", err)
		}
	})
}

func TestUserOldService_Login(t *testing.T) {
	repo := NewMockUserRepo()
	svc := service.NewUserOldService(repo)

	// Setup user
	email := "login_old@example.com"
	rawPassword := "secret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	repo.Create(&domain.User{
		Email:    email,
		Password: string(hashed),
	})

	t.Run("Success", func(t *testing.T) {
		tokens, err := svc.Login(email, rawPassword)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tokens.AccessToken == "" || tokens.RefreshToken == "" {
			t.Error("expected both access and refresh tokens")
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, err := svc.Login(email, "wrongpass")
		if err == nil {
			t.Error("expected error for invalid password")
		}
		if err.Error() != "invalid credentials" {
			t.Errorf("expected 'invalid credentials', got '%v'", err)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		_, err := svc.Login("nobody_old@example.com", "pass")
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

func TestUserOldService_RefreshToken(t *testing.T) {
	repo := NewMockUserRepo()
	svc := service.NewUserOldService(repo)

	email := "refresh_old@example.com"
	user := &domain.User{Email: email, Password: "hash"}
	repo.Create(user)

	// Helper to manually create valid token signed with "secret"
	secret := "secret"
	createToken := func(userID string, tokenType string, exp time.Duration) string {
		claims := jwt.MapClaims{
			"sub":  userID,
			"type": tokenType,
			"exp":  time.Now().Add(exp).Unix(),
		}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, _ := t.SignedString([]byte(secret))
		return s
	}

	t.Run("Success", func(t *testing.T) {
		validRefresh := createToken(user.ID.String(), "refresh", time.Hour)
		pair, err := svc.RefreshToken(validRefresh)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if pair.AccessToken == "" {
			t.Error("expected new access token")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := svc.RefreshToken("invalid.token.string")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})

	t.Run("Wrong Type (Access Token as Refresh)", func(t *testing.T) {
		accessToken := createToken(user.ID.String(), "access", time.Minute)
		_, err := svc.RefreshToken(accessToken)
		if err == nil {
			t.Error("expected error when using access token as refresh token")
		}
		if err.Error() != "invalid token type" {
			t.Errorf("expected 'invalid token type', got '%v'", err)
		}
	})
}
