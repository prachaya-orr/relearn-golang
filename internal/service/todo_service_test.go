package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"github.com/prachaya-orr/relearn-golang/internal/service"
)

// MockTodoRepository is a manual mock for testing
type MockTodoRepository struct {
	todos map[uuid.UUID]domain.Todo
}

func NewMockTodoRepo() *MockTodoRepository {
	return &MockTodoRepository{
		todos: make(map[uuid.UUID]domain.Todo),
	}
}

func (m *MockTodoRepository) Create(todo *domain.Todo) error {
	todo.ID = uuid.New()
	m.todos[todo.ID] = *todo
	return nil
}

func (m *MockTodoRepository) FindAll() ([]domain.Todo, error) {
	var list []domain.Todo
	for _, t := range m.todos {
		list = append(list, t)
	}
	return list, nil
}

func (m *MockTodoRepository) FindByID(id uuid.UUID) (*domain.Todo, error) {
	t, ok := m.todos[id]
	if !ok {
		return nil, nil // Not found
	}
	return &t, nil
}

func (m *MockTodoRepository) Update(todo *domain.Todo) error {
	m.todos[todo.ID] = *todo
	return nil
}

func (m *MockTodoRepository) Delete(id uuid.UUID) error {
	delete(m.todos, id)
	return nil
}

func (m *MockTodoRepository) DeleteAll() error {
	m.todos = make(map[uuid.UUID]domain.Todo)
	return nil
}

func TestCreateTodo(t *testing.T) {
	repo := NewMockTodoRepo()
	svc := service.NewTodoService(repo)

	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		todo, err := svc.Create("Test Todo", "Desc", userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if todo.ID == uuid.Nil {
			t.Error("expected ID to be set")
		}
		if todo.Title != "Test Todo" {
			t.Errorf("expected title 'Test Todo', got %s", todo.Title)
		}
	})

	t.Run("Empty Title", func(t *testing.T) {
		_, err := svc.Create("", "Desc", userID)
		if err == nil {
			t.Error("expected error for empty title")
		}
	})
}

func TestFindAll(t *testing.T) {
	repo := NewMockTodoRepo()
	svc := service.NewTodoService(repo)
	userID := uuid.New()

	// Seed data
	svc.Create("Todo 1", "Desc 1", userID)
	svc.Create("Todo 2", "Desc 2", userID)

	list, err := svc.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 todos, got %d", len(list))
	}
}

func TestFindByID(t *testing.T) {
	repo := NewMockTodoRepo()
	svc := service.NewTodoService(repo)
	userID := uuid.New()

	created, _ := svc.Create("Todo 1", "Desc 1", userID)

	t.Run("Found", func(t *testing.T) {
		found, err := svc.FindByID(created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if found == nil {
			t.Fatal("expected found todo, got nil")
		}
		if found.ID != created.ID {
			t.Errorf("expected ID %s, got %s", created.ID, found.ID)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		found, err := svc.FindByID(uuid.New())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if found != nil {
			t.Error("expected nil for non-existent ID")
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := NewMockTodoRepo()
	svc := service.NewTodoService(repo)
	userID := uuid.New()
	created, _ := svc.Create("Original", "Original Desc", userID)

	t.Run("Success", func(t *testing.T) {
		updated, err := svc.Update(created.ID, "Updated", "Updated Desc", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if updated.Title != "Updated" || updated.Description != "Updated Desc" || !updated.Completed {
			t.Error("expected fields to be updated")
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := svc.Update(uuid.New(), "Title", "Desc", true)
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
		if err.Error() != "todo not found" {
			t.Errorf("expected 'todo not found', got '%v'", err)
		}
	})
}

func TestDelete(t *testing.T) {
	repo := NewMockTodoRepo()
	svc := service.NewTodoService(repo)
	userID := uuid.New()
	created, _ := svc.Create("To Delete", "Desc", userID)

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it's gone
	found, _ := svc.FindByID(created.ID)
	if found != nil {
		t.Error("expected todo to be deleted")
	}
}

func TestDeleteAll(t *testing.T) {
	repo := NewMockTodoRepo()
	svc := service.NewTodoService(repo)
	userID := uuid.New()
	svc.Create("Todo 1", "Desc", userID)
	svc.Create("Todo 2", "Desc", userID)

	err := svc.DeleteAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	list, _ := svc.FindAll()
	if len(list) != 0 {
		t.Errorf("expected 0 todos after delete all, got %d", len(list))
	}
}
