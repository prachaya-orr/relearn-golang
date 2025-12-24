package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email    string    `gorm:"uniqueIndex;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserService interface {
	SignUp(email, password string) (*User, error)
	Login(email, password string) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
}
