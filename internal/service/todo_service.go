package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
)

// todoService implements domain.TodoService.
type todoService struct {
	repo domain.TodoRepository
}

// NewTodoService creates a new instance of TodoService.
func NewTodoService(repo domain.TodoRepository) domain.TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(title, description string, userID uuid.UUID) (*domain.Todo, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	todo := &domain.Todo{
		Title:       title,
		Description: description,
		UserID:      userID,
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) FindAll() ([]domain.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) FindByID(id uuid.UUID) (*domain.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) Update(id uuid.UUID, title, description string, completed bool) (*domain.Todo, error) {
	todo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}

	if title != "" {
		todo.Title = title
	}
	todo.Description = description
	todo.Completed = completed

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *todoService) DeleteAll() error {
	return s.repo.DeleteAll()
}
