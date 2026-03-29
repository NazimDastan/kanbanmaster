package middleware

import (
	"context"
	"net/http"
	"strings"

	"kanbanmaster/cmd/services"
)

func Auth(authService *services.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			// Try Authorization header first
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}

			// Fallback: query param for WebSocket connections
			if token == "" {
				token = r.URL.Query().Get("token")
			}

			if token == "" {
				http.Error(w, `{"error":"Authorization required","code":401}`, http.StatusUnauthorized)
				return
			}

			userID, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, `{"error":"Invalid or expired token","code":401}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
