package main

import (
	"fmt"
	"net/http"

	"hello-go-backend/internal/database"
	controllerRules "hello-go-backend/internal/rules/controller"
	controllerStatus "hello-go-backend/internal/status/controller"
	controllerUser "hello-go-backend/internal/user/controller"
	controllerUserStatus "hello-go-backend/internal/user_status/controller"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go Backend!")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := database.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

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
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/health", HealthHandler)

	//USERS
	userHandler := controllerUser.NewHandler(database.DB)
	http.HandleFunc("/user", userHandler.GetUserHandler)
	http.HandleFunc("/users", userHandler.GetAllUsersHandler)
	http.HandleFunc("/user/create", userHandler.CreateUserHandler)
	http.HandleFunc("/user/delete", userHandler.DeleteUserHandler)
	http.HandleFunc("/user/update", userHandler.UpdateUserHandler)
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
	http.HandleFunc("/user-status/update", userStatusHandler.Update)
	//USER STATUS
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
