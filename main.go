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

// corsMiddleware добавляет CORS-заголовки, чтобы index.html мог обращаться
// к API с другого origin (file://, http://localhost:5500 и т.п.) и отвечает
// на preflight-запросы (OPTIONS) без вызова основного обработчика.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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

	mux := http.NewServeMux()

	// Регистрируем обработчики
	// Служебные
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/health", HealthHandler)

	//USERS
	userHandler := controllerUser.NewHandler(database.DB)
	mux.HandleFunc("/user", userHandler.GetUserHandler)
	mux.HandleFunc("/users", userHandler.GetAllUsersHandler)
	mux.HandleFunc("/user/create", userHandler.CreateUserHandler)
	mux.HandleFunc("/user/delete", userHandler.DeleteUserHandler)
	mux.HandleFunc("/user/update", userHandler.UpdateUserHandler)
	// END USERS
	// STATUS
	statusHandler := controllerStatus.NewHandler(database.DB)
	mux.HandleFunc("/statuses", statusHandler.GetAll)
	mux.HandleFunc("/status", statusHandler.Get)
	mux.HandleFunc("/status/create", statusHandler.Create)
	mux.HandleFunc("/status/delete", statusHandler.Delete)
	mux.HandleFunc("/status/update", statusHandler.Update)
	// END STATUS
	//RULES
	rulesHandler := controllerRules.NewHandler(database.DB)
	mux.HandleFunc("/rules", rulesHandler.GetAll)
	mux.HandleFunc("/rule", rulesHandler.Get)
	mux.HandleFunc("/rule/create", rulesHandler.Create)
	mux.HandleFunc("/rule/delete", rulesHandler.Delete)
	mux.HandleFunc("/rule/update", rulesHandler.Update)
	//END RULES
	//USER STATUS
	userStatusHandler := controllerUserStatus.NewHandler(database.DB)
	mux.HandleFunc("/user-statuses", userStatusHandler.GetAll)
	mux.HandleFunc("/user-status", userStatusHandler.Get)
	mux.HandleFunc("/user-status/create", userStatusHandler.Create)
	mux.HandleFunc("/user-status/delete", userStatusHandler.Delete)
	mux.HandleFunc("/user-status/update", userStatusHandler.Update)
	//USER STATUS

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		panic(err)
	}
}
