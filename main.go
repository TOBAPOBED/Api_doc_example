package main

import (
	"log"
	"net/http"

	"api-doc-example/internal/handlers"
	_ "api-doc-example/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Todo API Service
// @version 1.0
// @description REST API для управления задачами (Todo) и пользователями, написанный на Go. Проект демонстрирует использование swaggo для автоматической генерации Swagger-документации.
// @BasePath /
// @schemes http https

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Users routes
	userHandler := handlers.NewUserHandler()
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", userHandler.ListUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Post("/", userHandler.CreateUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	// Todos routes
	todoHandler := handlers.NewTodoHandler()
	r.Route("/api/v1/todos", func(r chi.Router) {
		r.Get("/", todoHandler.ListTodos)
		r.Get("/{id}", todoHandler.GetTodo)
		r.Post("/", todoHandler.CreateTodo)
		r.Put("/{id}", todoHandler.UpdateTodo)
		r.Delete("/{id}", todoHandler.DeleteTodo)
	})

	// Swagger UI — работает с любым хостом (не привязан к localhost)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
	))

	// Раздача swagger.json — работает локально и на сервере
	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	log.Println("✅ Сервер запущен: http://localhost:8080")
	log.Println("✅ Swagger UI: http://localhost:8080/swagger/index.html")
	log.Println("✅ API JSON: http://localhost:8080/docs/swagger.json")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}