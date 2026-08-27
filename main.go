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
// @description Профессионально задокументированный REST API для управления задачами и пользователями.
// @host localhost:8080
// @BasePath /

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

	// Swagger UI с правильным маршрутом
	r.Get("/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
	))

	// Swagger JSON
	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	log.Println("✅ Сервер запущен: http://localhost:8080")
	log.Println("✅ Swagger UI: http://localhost:8080/swagger")
	log.Println("✅ API Docs JSON: http://localhost:8080/docs/swagger.json")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}