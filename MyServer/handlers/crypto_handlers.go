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

// EncryptRequest defines the request body for encrypting data
type EncryptRequest struct {
	KeyName string `json:"key_name"`
	Data    string `json:"data"`
}

// EncryptResponse defines the response body for encrypted data
type EncryptResponse struct {
	EncryptedData string `json:"encrypted_data"`
	Nonce         string `json:"nonce"`
}

// DecryptRequest defines the request body for decrypting data
type DecryptRequest struct {
	KeyName string `json:"key_name"`
	Data    string `json:"data"` // This is the encrypted data
	Nonce   string `json:"nonce"`
}

// DecryptResponse defines the response body for decrypted data
type DecryptResponse struct {
	DecryptedData string `json:"decrypted_data"`
}

// EncryptData handles the encryption of data using a specified key
func EncryptData(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	var req EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.KeyName == "" || req.Data == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Key name and data are required")
		return
	}

	key, err := database.GetKeyByName(req.KeyName, claims.UserID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if key == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Key not found or not owned by user")
		return
	}

	encryptedData, nonce, err := utils.Encrypt(key.KeyMaterial, []byte(req.Data))
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to encrypt data")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, EncryptResponse{
		EncryptedData: utils.EncodeToBase64(encryptedData),
		Nonce:         utils.EncodeToBase64(nonce),
	})
}

// DecryptData handles the decryption of data using a specified key
func DecryptData(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User claims not found")
		return
	}

	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.KeyName == "" || req.Data == "" || req.Nonce == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Key name, encrypted data, and nonce are required")
		return
	}

	key, err := database.GetKeyByName(req.KeyName, claims.UserID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if key == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Key not found or not owned by user")
		return
	}

	encryptedData, err := utils.DecodeFromBase64(req.Data)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid encrypted data format")
		return
	}

	nonce, err := utils.DecodeFromBase64(req.Nonce)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid nonce format")
		return
	}

	decryptedData, err := utils.Decrypt(key.KeyMaterial, encryptedData, nonce)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to decrypt data. Check key, data, and nonce.")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, DecryptResponse{DecryptedData: string(decryptedData)})
}

// RotateKey handles the rotation of an existing key for the authenticated user
func RotateKey(w http.ResponseWriter, r *http.Request) {
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

	// Generate new key material
	newKeyMaterial, err := utils.GenerateAESKey()
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to generate new key material")
		return
	}

	key.KeyMaterial = newKeyMaterial

	if err := database.UpdateKey(key); err != nil {
		if err == sql.ErrNoRows {
			middleware.RespondWithError(w, http.StatusNotFound, "Key not found for update or not owned by user")
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to rotate key")
		return
	}

	keyResp := database.KeyResponse{
		ID: key.ID,
		UserID: key.UserID,
		Name: key.Name,
		CreatedAt: key.CreatedAt,
	}
	middleware.RespondWithJSON(w, http.StatusOK, keyResp)
}