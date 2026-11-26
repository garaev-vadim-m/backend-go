package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"hello-go-backend/internal/database"
	"hello-go-backend/internal/middleware"
	controllerProduct "hello-go-backend/internal/product/controller"
	controllerRules "hello-go-backend/internal/rules/controller"
	"hello-go-backend/internal/session"
	controllerStatus "hello-go-backend/internal/status/controller"
	controllerUser "hello-go-backend/internal/user/controller"
	controllerUserStatus "hello-go-backend/internal/user_status/controller"

	controllerCart "hello-go-backend/internal/cart/controller"
	controllerOrder "hello-go-backend/internal/order/controller"
)

//go:embed web
var webFS embed.FS

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

func corsMiddleware(allowed map[string]bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Add("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseOrigins(s string) map[string]bool {
	m := map[string]bool{}
	for _, o := range strings.Split(s, ",") {
		if o = strings.TrimSpace(o); o != "" {
			m[o] = true
		}
	}
	return m
}

func getenvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func startSessionCleanup(db *sql.DB) {
	go func() {
		t := time.NewTicker(time.Hour)
		defer t.Stop()
		for range t.C {
			n, err := session.Cleanup(db, 30*24*time.Hour)
			if err != nil {
				log.Printf("session cleanup: %v", err)
				continue
			}
			if n > 0 {
				log.Printf("session cleanup: removed %d", n)
			}
		}
	}()
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

	startSessionCleanup(database.DB)

	if err := database.SeedAdmin(); err != nil {
		panic(err)
	}

	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	cookieSecure := getenvBool("COOKIE_SECURE", false)
	// Регистрируем обработчики
	// Служебные
	mux.HandleFunc("/hello", HomeHandler)
	mux.HandleFunc("/health", HealthHandler)
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	auth := middleware.NewAuth(database.DB, cookieSecure)
	//USERS
	userHandler := controllerUser.NewHandler(database.DB, cookieSecure)
	mux.Handle("GET /users", auth.RequireRole("admin", http.HandlerFunc(userHandler.GetAll)))
	mux.Handle("GET /user", auth.Require(http.HandlerFunc(userHandler.Get)))
	mux.Handle("GET /user/me", auth.Require(http.HandlerFunc(userHandler.Me)))
	mux.Handle("POST /user/create", auth.RequireRole("admin", http.HandlerFunc(userHandler.Create)))
	mux.Handle("DELETE /user/delete", auth.RequireRole("admin", http.HandlerFunc(userHandler.Delete)))
	mux.Handle("PUT /user/update", auth.RequireRole("admin", http.HandlerFunc(userHandler.Update)))
	// END USERS
	//AUTH
	mux.HandleFunc("POST /user/login", userHandler.Login)
	mux.HandleFunc("POST /user/logout", userHandler.Logout)
	mux.HandleFunc("POST /user/register", userHandler.Register)
	//End AUTH
	// STATUS
	statusHandler := controllerStatus.NewHandler(database.DB)
	mux.HandleFunc("GET /statuses", statusHandler.GetAll)
	mux.HandleFunc("GET /status", statusHandler.Get)
	mux.Handle("POST /status/create", auth.RequireRole("admin", http.HandlerFunc(statusHandler.Create)))
	mux.Handle("DELETE /status/delete", auth.RequireRole("admin", http.HandlerFunc(statusHandler.Delete)))
	mux.Handle("PUT /status/update", auth.RequireRole("admin", http.HandlerFunc(statusHandler.Update)))
	// END STATUS
	//RULES
	rulesHandler := controllerRules.NewHandler(database.DB)
	mux.HandleFunc("GET /rules", rulesHandler.GetAll)
	mux.HandleFunc("GET /rule", rulesHandler.Get)
	mux.Handle("POST /rule/create", auth.RequireRole("admin", http.HandlerFunc(rulesHandler.Create)))
	mux.Handle("DELETE /rule/delete", auth.RequireRole("admin", http.HandlerFunc(rulesHandler.Delete)))
	mux.Handle("PUT /rule/update", auth.RequireRole("admin", http.HandlerFunc(rulesHandler.Update)))
	//END RULES
	//USER STATUS
	userStatusHandler := controllerUserStatus.NewHandler(database.DB)
	mux.HandleFunc("GET /user-statuses", userStatusHandler.GetAll)
	mux.HandleFunc("GET /user-status", userStatusHandler.Get)
	mux.Handle("POST /user-status/create", auth.RequireRole("admin", http.HandlerFunc(userStatusHandler.Create)))
	mux.Handle("DELETE /user-status/delete", auth.RequireRole("admin", http.HandlerFunc(userStatusHandler.Delete)))
	mux.Handle("PUT /user-status/update", auth.RequireRole("admin", http.HandlerFunc(userStatusHandler.Update)))
	//END USER STATUS
	//PRODUCT
	productHandler := controllerProduct.NewHandler(database.DB)
	mux.HandleFunc("GET /products", productHandler.GetAll)
	mux.HandleFunc("GET /product", productHandler.Get)
	mux.Handle("POST /product/create", auth.RequireRole("admin", http.HandlerFunc(productHandler.Create)))
	mux.Handle("DELETE /product/delete", auth.RequireRole("admin", http.HandlerFunc(productHandler.Delete)))
	mux.Handle("PUT /product/update", auth.RequireRole("admin", http.HandlerFunc(productHandler.Update)))
	//END PRODUCT
	//CART
	cartHandler := controllerCart.NewHandler(database.DB)
	mux.Handle("GET /cart", auth.Require(http.HandlerFunc(cartHandler.Get)))
	mux.Handle("POST /cart/items", auth.Require(http.HandlerFunc(cartHandler.AddItem)))
	mux.Handle("DELETE /cart/items", auth.Require(http.HandlerFunc(cartHandler.RemoveItem)))
	//END CART
	//ORDER
	orderHandler := controllerOrder.NewHandler(database.DB)
	mux.Handle("POST /order/checkout", auth.Require(http.HandlerFunc(orderHandler.Checkout)))
	//END ORDER

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(parseOrigins(os.Getenv("CORS_ORIGINS")), mux)); err != nil {
		panic(err)
	}
}
