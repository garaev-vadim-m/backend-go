package main

import (
	"fmt"
	"net/http"

	"hello-go-backend/internal/database"
	controllerProduct "hello-go-backend/internal/product/controller"
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
	mux.HandleFunc("GET /user", userHandler.Get)
	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("POST /user/create", userHandler.Create)
	mux.HandleFunc("DELETE /user/delete", userHandler.Delete)
	mux.HandleFunc("PUT /user/update", userHandler.Update)
	// END USERS
	// STATUS
	statusHandler := controllerStatus.NewHandler(database.DB)
	mux.HandleFunc("GET /statuses", statusHandler.GetAll)
	mux.HandleFunc("GET /status", statusHandler.Get)
	mux.HandleFunc("POST /status/create", statusHandler.Create)
	mux.HandleFunc("DELETE /status/delete", statusHandler.Delete)
	mux.HandleFunc("PUT /status/update", statusHandler.Update)
	// END STATUS
	//RULES
	rulesHandler := controllerRules.NewHandler(database.DB)
	mux.HandleFunc("GET /rules", rulesHandler.GetAll)
	mux.HandleFunc("GET /rule", rulesHandler.Get)
	mux.HandleFunc("POST /rule/create", rulesHandler.Create)
	mux.HandleFunc("DELETE /rule/delete", rulesHandler.Delete)
	mux.HandleFunc("PUT /rule/update", rulesHandler.Update)
	//END RULES
	//USER STATUS
	userStatusHandler := controllerUserStatus.NewHandler(database.DB)
	mux.HandleFunc("GET /user-statuses", userStatusHandler.GetAll)
	mux.HandleFunc("GET /user-status", userStatusHandler.Get)
	mux.HandleFunc("POST /user-status/create", userStatusHandler.Create)
	mux.HandleFunc("DELETE /user-status/delete", userStatusHandler.Delete)
	mux.HandleFunc("PUT /user-status/update", userStatusHandler.Update)
	//END USER STATUS
	//PRODUCT
	productHandler := controllerProduct.NewHandler(database.DB)
	mux.HandleFunc("GET /products", productHandler.GetAll)
	mux.HandleFunc("GET /product", productHandler.Get)
	mux.HandleFunc("POST /product/create", productHandler.Create)
	mux.HandleFunc("DELETE /product/delete", productHandler.Delete)
	mux.HandleFunc("PUT /product/update", productHandler.Update)
	//END PRODUCT

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		panic(err)
	}
}
