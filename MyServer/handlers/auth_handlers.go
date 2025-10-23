package handlers 

import (
	"encoding/json"
	"net/http"

	"github.com/anurag/magicgate/MyServer/config"
	"github.com/anurag/magicgate/MyServer/database"
	"github.com/anurag/magicgate/MyServer/middleware"
	"github.com/anurag/magicgate/MyServer/utils"
)

// LoginRequest defines the request body for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse defines the response body for successful login
type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles user authentication and JWT generation
func Login(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if req.Username == "" || req.Password == "" {
			middleware.RespondWithError(w, http.StatusBadRequest, "Username and password are required")
			return
		}

		user, err := database.GetUserByUsername(req.Username)
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}
		if user == nil {
			middleware.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			middleware.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Username, []byte(cfg.JWTSecret))
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		middleware.RespondWithJSON(w, http.StatusOK, LoginResponse{Token: token})
	}
}