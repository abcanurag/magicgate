package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anurag/magicgate/MyServer/database"
	"github.com/anurag/magicgate/MyServer/middleware"
	"github.com/anurag/magicgate/MyServer/utils"
	"github.com/gorilla/mux"
)

// KeyCreateRequest defines the request body for creating a key
type KeyCreateRequest struct {
	Name string `json:"name"`
}

// KeyUpdateRequest defines the request body for updating a key
type KeyUpdateRequest struct {
	Name string `json:"name"`
}

// CreateKey handles the creation of a new cryptographic key for the authenticated user
func CreateKey(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	var req KeyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Name == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Key name is required")
		return
	}

	// Generate a new AES key
	keyMaterial, err := utils.GenerateAESKey()
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to generate key material")
		return
	}

	key := &database.Key{
		UserID:      claims.UserID,
		Name:        req.Name,
		KeyMaterial: keyMaterial,
	}

	if err := database.CreateKey(key); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create key")
		return
	}

	// Respond with KeyResponse to avoid exposing raw key material
	keyResp := database.KeyResponse{
		ID:        key.ID,
		UserID:    key.UserID,
		Name:      key.Name,
		CreatedAt: key.CreatedAt,
	}
	middleware.RespondWithJSON(w, http.StatusCreated, keyResp)
}

// GetKey handles retrieving a specific key for the authenticated user
func GetKey(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	key, err := database.GetKeyByID(id, claims.UserID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if key == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Key not found or not owned by user")
		return
	}

	keyResp := database.KeyResponse{
		ID:        key.ID,
		UserID:    key.UserID,
		Name:      key.Name,
		CreatedAt: key.CreatedAt,
	}
	middleware.RespondWithJSON(w, http.StatusOK, keyResp)
}

// GetAllKeys handles retrieving all keys for the authenticated user
func GetAllKeys(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	keys, err := database.GetAllKeysForUser(claims.UserID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	keyResponses := make([]database.KeyResponse, len(keys))
	for i, key := range keys {
		keyResponses[i] = database.KeyResponse{
			ID:        key.ID,
			UserID:    key.UserID,
			Name:      key.Name,
			CreatedAt: key.CreatedAt,
		}
	}
	middleware.RespondWithJSON(w, http.StatusOK, keyResponses)
}

// UpdateKey handles updating a key's name for the authenticated user
func UpdateKey(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	var req KeyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Name == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Key name is required for update")
		return
	}

	key, err := database.GetKeyByID(id, claims.UserID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if key == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Key not found or not owned by user")
		return
	}

	key.Name = req.Name
	// Note: KeyMaterial is not updated via this endpoint for simplicity.
	// A separate endpoint or flow might be needed for key rotation/regeneration.

	if err := database.UpdateKey(key); err != nil {
		if err == sql.ErrNoRows {
			middleware.RespondWithError(w, http.StatusNotFound, "Key not found for update or not owned by user")
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update key")
		return
	}

	keyResp := database.KeyResponse{
		ID:        key.ID,
		UserID:    key.UserID,
		Name:      key.Name,
		CreatedAt: key.CreatedAt,
	}
	middleware.RespondWithJSON(w, http.StatusOK, keyResp)
}

// DeleteKey handles deleting a key for the authenticated user
func DeleteKey(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	err = database.DeleteKey(id, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.RespondWithError(w, http.StatusNotFound, "Key not found or not owned by user")
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete key")
		return
	}

	middleware.RespondWithJSON(w, http.StatusNoContent, nil)
}