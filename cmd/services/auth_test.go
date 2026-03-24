package services

import (
	"testing"
	"time"

	"kanbanmaster/cmd/config"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAndParseToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := &AuthService{cfg: cfg}

	userID := "test-user-id"
	tokenStr, err := svc.createToken(userID, 15*time.Minute)
	if err != nil {
		t.Fatalf("createToken failed: %v", err)
	}

	if tokenStr == "" {
		t.Fatal("token should not be empty")
	}

	claims, err := svc.parseToken(tokenStr)
	if err != nil {
		t.Fatalf("parseToken failed: %v", err)
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub != userID {
		t.Fatalf("expected sub=%s, got %v", userID, claims["sub"])
	}
}

func TestParseInvalidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := &AuthService{cfg: cfg}

	_, err := svc.parseToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestParseExpiredToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := &AuthService{cfg: cfg}

	claims := jwt.MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(cfg.JWTSecret))

	_, err := svc.parseToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := &AuthService{cfg: cfg}

	tokenStr, _ := svc.createToken("user-123", 15*time.Minute)

	userID, err := svc.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if userID != "user-123" {
		t.Fatalf("expected user-123, got %s", userID)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := &AuthService{cfg: cfg}

	_, err := svc.ValidateToken("bad-token")
	if err == nil {
		t.Fatal("expected error for bad token")
	}
}
