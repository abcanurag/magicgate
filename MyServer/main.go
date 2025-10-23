package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anurag/magicgate/MyServer/config"
	"github.com/anurag/magicgate/MyServer/database"
	"github.com/anurag/magicgate/MyServer/handlers"
	"github.com/anurag/magicgate/MyServer/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(cfg.DatabaseURL)
	defer database.CloseDB()

	// Setup router
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/register", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/login", handlers.Login(cfg)).Methods("POST")

	// Authenticated routes
	authRouter := r.PathPrefix("/api").Subrouter()
	authRouter.Use(middleware.AuthMiddleware(cfg))

	// User CRUD (authenticated, but for simplicity, any authenticated user can access any user ID for now)
	// A more secure implementation would check if the ID in the path matches the authenticated user's ID
	// or if the user has an 'admin' role.
	authRouter.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	authRouter.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	authRouter.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	authRouter.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	// Key CRUD (authenticated and user-specific)
	authRouter.HandleFunc("/keys", handlers.CreateKey).Methods("POST")
	authRouter.HandleFunc("/keys", handlers.GetAllKeys).Methods("GET")
	authRouter.HandleFunc("/keys/{id}", handlers.GetKey).Methods("GET")
	authRouter.HandleFunc("/keys/{id}", handlers.UpdateKey).Methods("PUT")
	authRouter.HandleFunc("/keys/{id}", handlers.DeleteKey).Methods("DELETE")

	// Crypto operations (authenticated and user-specific)
	authRouter.HandleFunc("/encrypt", handlers.Encrypt(cfg)).Methods("POST")
	authRouter.HandleFunc("/decrypt", handlers.Decrypt(cfg)).Methods("POST")

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}