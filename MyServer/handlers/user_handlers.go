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

// UserCreateRequest defines the request body for creating a user
type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserUpdateRequest defines the request body for updating a user
type UserUpdateRequest struct {
	Username string `json:"username"`
}

// CreateUser handles the creation of a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Password == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Check if user already exists
	existingUser, err := database.GetUserByUsername(req.Username)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error checking user existence")
		return
	}
	if existingUser != nil {
		middleware.RespondWithError(w, http.StatusConflict, "User with this username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := &database.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	if err := database.CreateUser(user); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, user)
}

// GetUser handles retrieving a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := database.GetUserByID(id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}

// GetAllUsers handles retrieving all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	middleware.RespondWithJSON(w, http.StatusOK, users)
}

// UpdateUser handles updating a user's information
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" {
		middleware.RespondWithError(w, http.StatusBadRequest, "Username is required for update")
		return
	}

	user, err := database.GetUserByID(id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	user.Username = req.Username
	if err := database.UpdateUser(user); err != nil {
		if err == sql.ErrNoRows {
			middleware.RespondWithError(w, http.StatusNotFound, "User not found for update")
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles deleting a user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = database.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	middleware.RespondWithJSON(w, http.StatusNoContent, nil)
}