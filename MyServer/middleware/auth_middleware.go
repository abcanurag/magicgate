package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anurag/magicgate/MyServer/config"
	"github.com/anurag/magicgate/MyServer/utils"
)

// UserContextKey is a custom type for context keys to avoid collisions
type UserContextKey string

const (
	// AuthenticatedUserKey is the key used to store the authenticated user's claims in the context
	AuthenticatedUserKey UserContextKey = "authenticatedUser"
)

// AuthMiddleware validates JWT tokens and adds user claims to the request context
func AuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			return
		}

		// Add user claims to the request context
		ctx := context.WithValue(r.Context(), AuthenticatedUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserClaimsFromContext retrieves user claims from the request context
func GetUserClaimsFromContext(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(AuthenticatedUserKey).(*utils.Claims)
	return claims, ok
}

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}