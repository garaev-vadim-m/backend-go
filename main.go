package main

import (
	"fmt"
	"net/http"

	"hello-go-backend/internal/database"
	"hello-go-backend/internal/handlers"
	controllerRules "hello-go-backend/internal/rules/controller"
	controllerStatus "hello-go-backend/internal/status/controller"
	controllerUserStatus "hello-go-backend/internal/user_status/controller"
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
	http.HandleFunc("/user", handlers.GetUserHandler)
	http.HandleFunc("/users", handlers.GetAllUsersHandler)
	http.HandleFunc("/user/create", handlers.CreateUserHandler)
	http.HandleFunc("/user/delete", handlers.DeleteUserHandler)
	http.HandleFunc("/user/update", handlers.UpdateUserHandler)
	// END USERS
	// STATUS
	statusHandler := controllerStatus.NewHandler(database.DB)
	http.HandleFunc("/statuses", statusHandler.GetAll)
	http.HandleFunc("/status", statusHandler.Get)
	http.HandleFunc("/status/create", statusHandler.Create)
	http.HandleFunc("/status/delete", statusHandler.Delete)
	http.HandleFunc("/status/update", statusHandler.Update)
	// END STATUS
	//RULES
	rulesHandler := controllerRules.NewHandler(database.DB)
	http.HandleFunc("/rules", rulesHandler.GetAll)
	http.HandleFunc("/rule", rulesHandler.Get)
	http.HandleFunc("/rule/create", rulesHandler.Create)
	http.HandleFunc("/rule/delete", rulesHandler.Delete)
	http.HandleFunc("/rule/update", rulesHandler.Update)
	//END RULES
	//USER STATUS
	userStatusHandler := controllerUserStatus.NewHandler(database.DB)
	http.HandleFunc("/user-statuses", userStatusHandler.GetAll)
	http.HandleFunc("/user-status", userStatusHandler.Get)
	http.HandleFunc("/user-status/create", userStatusHandler.Create)
	http.HandleFunc("/user-status/delete", userStatusHandler.Delete)
	//USER STATUS
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
