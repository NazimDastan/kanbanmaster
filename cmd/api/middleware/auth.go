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
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header required","code":401}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"Invalid authorization format","code":401}`, http.StatusUnauthorized)
				return
			}

			userID, err := authService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"Invalid or expired token","code":401}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
