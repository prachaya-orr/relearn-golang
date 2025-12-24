package repository

import (
	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new GORM repository.
func NewTodoRepository(db *gorm.DB) domain.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *domain.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) FindAll() ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

func (r *todoRepository) FindByID(id uuid.UUID) (*domain.Todo, error) {
	var todo domain.Todo
	err := r.db.First(&todo, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Update(todo *domain.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Todo{}, "id = ?", id).Error
}

func (r *todoRepository) DeleteAll() error {
	return r.db.Exec("DELETE FROM todos").Error
}
