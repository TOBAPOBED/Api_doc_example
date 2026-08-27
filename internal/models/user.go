package models

import "time"

// User представляет собой пользователя системы.
// @Description Модель данных пользователя.
type User struct {
	ID        int       `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"Иван Иванов"`
	CreatedAt time.Time `json:"created_at" example:"2023-10-25T10:00:00Z"`
}

// CreateUserRequest запрос на создание пользователя.
// @Description Данные для регистрации нового пользователя.
type CreateUserRequest struct {
	Email string `json:"email" example:"newuser@example.com" validate:"required,email"`
	Name  string `json:"name" example:"Петр Петров" validate:"required,max=100"`
}

// UpdateUserRequest запрос на обновление пользователя.
// @Description Поля для частичного обновления данных пользователя.
type UpdateUserRequest struct {
	Email *string `json:"email" example:"updated@example.com" validate:"omitempty,email"`
	Name  *string `json:"name" example:"Новое Имя" validate:"omitempty,max=100"`
}

// ErrorResponse стандартный формат ответа при ошибке.
// @Description Структура ответа сервера при возникновении ошибки.
type ErrorResponse struct {
	Error   string `json:"error" example:"bad_request"`
	Message string `json:"message" example:"Неверный ID"`
}

// SuccessResponse стандартный формат успешного ответа.
// @Description Универсальная обертка для успешных ответов API.
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}