package database

import (
	"database/sql"
	"fmt"
)

// CreateUser inserts a new user into the database
func CreateUser(user *User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := DB.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password_hash, created_at FROM users WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]User, error) {
	rows, err := DB.Query(`SELECT id, username, password_hash, created_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates an existing user's username (password update would be separate)
func UpdateUser(user *User) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	result, err := DB.Exec(query, user.Username, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // User not found for update
	}
	return nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // User not found for delete
	}
	return nil
}