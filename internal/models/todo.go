package models

import "time"

// Todo представляет собой задачу в системе управления задачами.
// @Description Модель данных задачи, содержащая всю информацию о таске.
type Todo struct {
	ID          int       `json:"id" example:"1"`
	UserID      int       `json:"user_id" example:"10"`
	Title       string    `json:"title" example:"Изучить Swagger в Go"`
	Description string    `json:"description" example:"Прочитать документацию swaggo"`
	Done        bool      `json:"done" example:"false"`
	CreatedAt   time.Time `json:"created_at" example:"2023-10-25T10:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-10-25T10:00:00Z"`
}

// CreateTodoRequest запрос на создание новой задачи.
// @Description Данные, необходимые для создания задачи.
type CreateTodoRequest struct {
	UserID      int    `json:"user_id" example:"10" validate:"required"`
	Title       string `json:"title" example:"Новая задача" validate:"required,max=200"`
	Description string `json:"description" example:"Описание задачи" validate:"max=1000"`
	Done        bool   `json:"done" example:"false"`
}

// UpdateTodoRequest запрос на обновление задачи.
// @Description Поля для частичного обновления задачи (все поля опциональны).
type UpdateTodoRequest struct {
	Title       *string `json:"title" example:"Обновленный заголовок" validate:"omitempty,max=200"`
	Description *string `json:"description" example:"Обновленное описание" validate:"omitempty,max=1000"`
	Done        *bool   `json:"done" example:"true"`
}