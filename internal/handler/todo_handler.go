package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
)

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required" example:"Buy milk"`
	Description string `json:"description" example:"Go to the store"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title" example:"Buy almond milk"`
	Description string `json:"description" example:"Go to the organic store"`
	Completed   bool   `json:"completed" example:"true"`
}

type TodoHandler struct {
	svc domain.TodoService
}

// NewTodoHandler creates a new TodoHandler.
func NewTodoHandler(svc domain.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// Create handles POST /todos
// @Summary Create a new todo
// @Description Create a new todo with the input payload
// @Tags todos
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param todo body CreateTodoRequest true "Create Todo"
// @Success 201 {object} domain.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	var req CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	todo, err := h.svc.Create(req.Title, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// FindAll handles GET /todos
// @Summary List all todos
// @Description Get all todos
// @Tags todos
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} domain.Todo
// @Failure 500 {object} map[string]string
// @Router /todos [get]
func (h *TodoHandler) FindAll(c *gin.Context) {
	todos, err := h.svc.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// FindByID handles GET /todos/:id
// @Summary Get a todo
// @Description Get a todo by ID
// @Tags todos
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} domain.Todo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [get]
func (h *TodoHandler) FindByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	todo, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if todo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Update handles PUT /todos/:id
// @Summary Update a todo
// @Description Update a todo by ID
// @Tags todos
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param todo body UpdateTodoRequest true "Update Todo"
// @Success 200 {object} domain.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req UpdateTodoRequest

	// Note: For partial updates, strictly you might want PATCH and pointers,
	// but for simplicity in this boiler plate we accept zero values.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.svc.Update(id, req.Title, req.Description, req.Completed)
	if err != nil {
		// In a real app we would check for specific errors like 'not found'
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Delete handles DELETE /todos/:id
// @Summary Delete a todo
// @Description Delete a todo by ID
// @Tags todos
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAll handles DELETE /todos
// @Summary Delete all todos
// @Description Delete all todos in the database (Requires API Key)
// @Tags todos
// @Produce  json
// @Security BearerAuth
// @Param X-API-KEY header string true "API Key"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [delete]
func (h *TodoHandler) DeleteAll(c *gin.Context) {
	// API Key Authentication
	apiKey := c.GetHeader("X-API-KEY")
	if apiKey == "" {
		apiKey = c.Query("api-key") // fallback to query param if needed
	}

	if apiKey != "delete" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API Key"})
		return
	}

	if err := h.svc.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
