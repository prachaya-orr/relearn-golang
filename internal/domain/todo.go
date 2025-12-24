package domain

import "github.com/google/uuid"

// Todo represents a task in the system.
type Todo struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
}

// TodoRepository defines the interface for database operations.
type TodoRepository interface {
	Create(todo *Todo) error
	FindAll() ([]Todo, error)
	FindByID(id uuid.UUID) (*Todo, error)
	Update(todo *Todo) error
	Delete(id uuid.UUID) error
	DeleteAll() error
}

// TodoService defines the interface for business logic.
type TodoService interface {
	Create(title, description string, userID uuid.UUID) (*Todo, error)
	FindAll() ([]Todo, error)
	FindByID(id uuid.UUID) (*Todo, error)
	Update(id uuid.UUID, title, description string, completed bool) (*Todo, error)
	Delete(id uuid.UUID) error
	DeleteAll() error
}
