package main

import (
	"fmt"
	"net/http"

	"hello-go-backend/internal/database"
	"hello-go-backend/internal/handlers"
	"hello-go-backend/internal/status"
)

func main() {
	// Инициализируем БД
	if err := database.InitDB(); err != nil {
		panic(err)
	}
	defer database.CloseDB()

	// Запускаем миграции
	if err := database.RunMigrations(); err != nil {
		panic(err)
	}

	// Регистрируем обработчики
	// Служебные
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	//USERS
	// Get
	http.HandleFunc("/user", handlers.GetUserHandler)
	http.HandleFunc("/users", handlers.GetAllUsersHandler)
	// Post
	http.HandleFunc("/user/create", handlers.CreateUserHandler)
	//Delete
	http.HandleFunc("/user/delete", handlers.DeleteUserHandler)
	// Put
	http.HandleFunc("/user/update", handlers.UpdateUserHandler)
	// END USERS
	statusHandler := status.NewHandler(database.DB)
	// STATUS
	http.HandleFunc("/statuses", statusHandler.StatusesHandler)
	http.HandleFunc("/status", statusHandler.StatusHandler)
	http.HandleFunc("/status/create", statusHandler.CreateStatusHandler)
	http.HandleFunc("/status/delete", statusHandler.DeleteStatusHandler)
	http.HandleFunc("/status/update", statusHandler.UpdateStatusHandler)
	// END STATUS
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
