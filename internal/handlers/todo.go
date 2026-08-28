package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"api-doc-example/internal/models"
)

type TodoHandler struct {
	mu     sync.RWMutex
	todos  map[int]models.Todo
	nextID int
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todos:  make(map[int]models.Todo),
		nextID: 1,
	}
}

// ListTodos возвращает список всех задач
// @Summary Получить список всех задач
// @Description Возвращает массив всех существующих задач в системе
// @Tags Todos
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse "Список задач успешно получен"
// @Router /api/v1/todos [get]
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	var todos []models.Todo
	for _, todo := range h.todos {
		todos = append(todos, todo)
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todos,
	})
}

// GetTodo возвращает задачу по её идентификатору
// @Summary Получить задачу по ID
// @Description Возвращает детальную информацию о задаче по её уникальному ID
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path int true "Уникальный идентификатор задачи"
// @Success 200 {object} models.SuccessResponse "Задача успешно найдена"
// @Failure 400 {object} models.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} models.ErrorResponse "Задача не найдена"
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Неверный формат ID",
		})
		return
	}

	h.mu.RLock()
	todo, exists := h.todos[id]
	h.mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Задача не найдена",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// CreateTodo создает новую задачу
// @Summary Создать новую задачу
// @Description Добавляет новую задачу в хранилище и возвращает её с присвоенным ID
// @Tags Todos
// @Accept json
// @Produce json
// @Param request body models.CreateTodoRequest true "Данные для создания задачи"
// @Success 201 {object} models.SuccessResponse "Задача успешно создана"
// @Failure 400 {object} models.ErrorResponse "Невалидный запрос"
// @Router /api/v1/todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Невалидный JSON",
		})
		return
	}

	// Валидация
	if req.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Title обязателен",
		})
		return
	}
	if len(req.Title) > 200 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Title не должен превышать 200 символов",
		})
		return
	}
	if len(req.Description) > 1000 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Description не должен превышать 1000 символов",
		})
		return
	}

	h.mu.Lock()
	todo := models.Todo{
		ID:          h.nextID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	h.todos[h.nextID] = todo
	h.nextID++
	h.mu.Unlock()

	log.Printf("Создана новая задача: %+v", todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// UpdateTodo обновляет существующую задачу
// @Summary Обновить задачу по ID
// @Description Частично обновляет данные задачи (поддерживаются title, description, done)
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path int true "Уникальный идентификатор задачи"
// @Param request body models.UpdateTodoRequest true "Данные для обновления задачи"
// @Success 200 {object} models.SuccessResponse "Задача успешно обновлена"
// @Failure 400 {object} models.ErrorResponse "Неверный формат ID или данных"
// @Failure 404 {object} models.ErrorResponse "Задача не найдена"
// @Router /api/v1/todos/{id} [put]
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Неверный формат ID",
		})
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Невалидный JSON",
		})
		return
	}

	// Валидация
	if req.Title != nil && len(*req.Title) > 200 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Title не должен превышать 200 символов",
		})
		return
	}
	if req.Description != nil && len(*req.Description) > 1000 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Description не должен превышать 1000 символов",
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	todo, exists := h.todos[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Задача не найдена",
		})
		return
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Done != nil {
		todo.Done = *req.Done
	}
	todo.UpdatedAt = time.Now()

	h.todos[id] = todo

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// DeleteTodo удаляет задачу
// @Summary Удалить задачу по ID
// @Description Безвозвратно удаляет задачу из хранилища
// @Tags Todos
// @Param id path int true "Уникальный идентификатор задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} models.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} models.ErrorResponse "Задача не найдена"
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Неверный формат ID",
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.todos[id]; !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Задача не найдена",
		})
		return
	}

	delete(h.todos, id)

	w.WriteHeader(http.StatusNoContent)
}