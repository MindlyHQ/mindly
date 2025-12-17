package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mindly/api/internal/models"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Получаем контекст из запроса
	ctx := r.Context()

	// Разбираем запрос
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидируем данные
	if err := validateRegisterRequest(req); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)"

	err := h.DB.QueryRowContext(ctx, query, strings.ToLower(req.Email), req.Username).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking user existence: %v", err)
		sendJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
		sendJSONError(w, "User with this email or username already exists", http.StatusConflict)
		return
	}

	// Хешируем пароль
	passwordHash, err := models.HashPassword(req.Password)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		sendJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базе данных
	var id string
	var createdAt, updatedAt time.Time

	insertQuery := `
		INSERT INTO users (
			email, 
			username, 
			password_hash, 
			full_name, 
			score, 
			current_streak, 
			best_streak
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err = h.DB.QueryRowContext(ctx, insertQuery,
		strings.ToLower(req.Email),
		req.Username,
		passwordHash,
		req.FullName,
		0,
		0,
		0,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		log.Printf("Database insert error: %v", err)
		sendJSONError(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем объект пользователя для ответа
	user := models.User{
		ID:            id,
		Email:         req.Email,
		Username:      req.Username,
		PasswordHash:  passwordHash,
		FullName:      req.FullName,
		Score:         0,
		CurrentStreak: 0,
		BestStreak:    0,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	// Отправляем успешный ответ
	response := models.APIResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    user,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func validateRegisterRequest(req models.RegisterRequest) error {
	// Проверяем email
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return fmt.Errorf("invalid email format")
	}

	// Проверяем username
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(req.Username) > 50 {
		return fmt.Errorf("username must be less than 50 characters")
	}

	// Проверяем пароль
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	return nil
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	response := models.APIResponse{
		Status: "error",
		Error:  message,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
