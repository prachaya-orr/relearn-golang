package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"github.com/prachaya-orr/relearn-golang/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a manual mock for testing UserService
type MockUserRepository struct {
	users map[string]*domain.User // Email -> User
}

func NewMockUserRepo() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("email already registered")
	}
	user.ID = uuid.New()
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	// Return a copy to simulate retrieval
	u := *user
	return &u, nil
}

func TestUserService_SignUp(t *testing.T) {
	repo := NewMockUserRepo()
	svc := service.NewUserService(repo)

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
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
		// Already created in previous step? No, repo is shared?
		// Actually, I should probably reset repo or handle it.
		// NewMockUserRepo() is called at top level, checking if state persists.
		// Ah, for "Duplicate Email", I need to ensure one exists first.

		// Let's create a fresh one for this specific sub-test to be clean slightly or reuse?
		// reusing 'repo' from outer scope.

		email := "dup@example.com"
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

func TestUserService_Login(t *testing.T) {
	repo := NewMockUserRepo()
	svc := service.NewUserService(repo)

	// Setup user
	email := "login@example.com"
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
		_, err := svc.Login("nobody@example.com", "pass")
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

func TestUserService_RefreshToken(t *testing.T) {
	repo := NewMockUserRepo()
	svc := service.NewUserService(repo)

	// We need a valid refresh token first. Login gives us one.
	email := "refresh@example.com"
	repo.Create(&domain.User{
		Email:    email,
		Password: "hashedpassword", // Doesn't matter for refresh if we bypass login, but login helper uses it.
		// Actually, let's just manually create a token or use Login helper from service if possible.
		// But repo needs to contain the user because token claims have sub=UserID,
		// and RefreshToken logic parses token -> gets ID -> genToken(ID).
		// Wait, newRefreshTokenFunc checks parsing.
		// It does NOT seem to check if user still exists in repo in the implementation I read?
		// Let's check user_service.go again.
		// ...
		// func newRefreshTokenFunc(...) {
		//    ... parse ... verify claims ...
		//    userID, err := uuid.Parse(sub)
		//    return genToken(userID)
		// }
		// It does NOT check DB for user existence! Just generates new token for that ID.
	})
	// But we need the ID to match what's in the token.
	// Let's create a user properly via repo to get an ID?
	// Or just manually construct the user.
	user := &domain.User{Email: email, Password: "hash"}
	repo.Create(user) // assigns ID

	// We can't easily generate a valid signed token without the secret used by the service.
	// The service uses os.Getenv("JWT_SECRET") or "secret".
	// So we can replicate that generation here.

	secret := "secret" // default in service

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
