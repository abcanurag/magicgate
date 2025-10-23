package database

import (
	"database/sql"
	"fmt"
)

// CreateKey inserts a new cryptographic key into the database
func CreateKey(key *Key) error {
	query := `INSERT INTO keys (user_id, name, key_material) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := DB.QueryRow(query, key.UserID, key.Name, key.KeyMaterial).Scan(&key.ID, &key.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}
	return nil
}

// GetKeyByID retrieves a key by its ID and user ID
func GetKeyByID(id, userID int) (*Key, error) {
	key := &Key{}
	query := `SELECT id, user_id, name, key_material, created_at FROM keys WHERE id = $1 AND user_id = $2`
	err := DB.QueryRow(query, id, userID).Scan(&key.ID, &key.UserID, &key.Name, &key.KeyMaterial, &key.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Key not found for this user
		}
		return nil, fmt.Errorf("failed to get key by ID: %w", err)
	}
	return key, nil
}

// GetKeyByName retrieves a key by its name and user ID
func GetKeyByName(name string, userID int) (*Key, error) {
	key := &Key{}
	query := `SELECT id, user_id, name, key_material, created_at FROM keys WHERE name = $1 AND user_id = $2`
	err := DB.QueryRow(query, name, userID).Scan(&key.ID, &key.UserID, &key.Name, &key.KeyMaterial, &key.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Key not found for this user
		}
		return nil, fmt.Errorf("failed to get key by name: %w", err)
	}
	return key, nil
}

// GetAllKeysForUser retrieves all keys for a specific user
func GetAllKeysForUser(userID int) ([]Key, error) {
	rows, err := DB.Query(`SELECT id, user_id, name, key_material, created_at FROM keys WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all keys for user: %w", err)
	}
	defer rows.Close()

	keys := []Key{}
	for rows.Next() {
		key := Key{}
		if err := rows.Scan(&key.ID, &key.UserID, &key.Name, &key.KeyMaterial, &key.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan key row: %w", err)
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// UpdateKey updates an existing key's name or material
func UpdateKey(key *Key) error {
	query := `UPDATE keys SET name = $1, key_material = $2 WHERE id = $3 AND user_id = $4`
	result, err := DB.Exec(query, key.Name, key.KeyMaterial, key.ID, key.UserID)
	if err != nil {
		return fmt.Errorf("failed to update key: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Key not found for update
	}
	return nil
}

// DeleteKey deletes a key by its ID and user ID
func DeleteKey(id, userID int) error {
	query := `DELETE FROM keys WHERE id = $1 AND user_id = $2`
	result, err := DB.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Key not found for delete
	}
	return nil
}