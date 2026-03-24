package services

import (
	"database/sql"
	"errors"
	"time"

	"kanbanmaster/cmd/config"
	"kanbanmaster/cmd/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidToken    = errors.New("invalid token")
)

type AuthService struct {
	db  *sql.DB
	cfg *config.Config
}

func NewAuthService(db *sql.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         models.User  `json:"user"`
}

func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", input.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.QueryRow(
		`INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3)
		 RETURNING id, email, password_hash, name, avatar_url, created_at, updated_at`,
		input.Email, string(hash), input.Name,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(user)
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	var user models.User
	err := s.db.QueryRow(
		"SELECT id, email, password_hash, name, avatar_url, created_at, updated_at FROM users WHERE email = $1",
		input.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return s.generateTokens(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	var user models.User
	err = s.db.QueryRow(
		"SELECT id, email, password_hash, name, avatar_url, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokens(user)
}

func (s *AuthService) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		"SELECT id, email, password_hash, name, avatar_url, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateAvatar(userID, avatarURL string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2
		 RETURNING id, email, password_hash, name, avatar_url, created_at, updated_at`,
		avatarURL, userID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateProfile(userID, name, email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3
		 RETURNING id, email, password_hash, name, avatar_url, created_at, updated_at`,
		name, email, userID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) ChangePassword(userID, currentPassword, newPassword string) error {
	var hash string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&hash)
	if err != nil {
		return ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(currentPassword)); err != nil {
		return ErrInvalidPassword
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2", string(newHash), userID)
	return err
}

func (s *AuthService) ValidateToken(tokenStr string) (string, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return "", ErrInvalidToken
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", ErrInvalidToken
	}
	return userID, nil
}

func (s *AuthService) generateTokens(user models.User) (*AuthResponse, error) {
	accessToken, err := s.createToken(user.ID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createToken(user.ID, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) createToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
