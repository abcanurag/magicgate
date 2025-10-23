package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB holds the database connection pool
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(databaseURL string) {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	createTables()
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}

// createTables creates necessary tables if they don't exist
func createTables() {
	userTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	keyTableSQL := `
	CREATE TABLE IF NOT EXISTS keys (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		name VARCHAR(255) NOT NULL,
		key_material BYTEA NOT NULL, -- Storing raw key material (e.g., AES key)
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, name) -- A user cannot have two keys with the same name
	);`

	_, err := DB.Exec(userTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("Users table checked/created.")

	_, err = DB.Exec(keyTableSQL)
	if err != nil {
		log.Fatalf("Error creating keys table: %v", err)
	}
	log.Println("Keys table checked/created.")
}