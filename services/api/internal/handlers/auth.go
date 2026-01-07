package handlers

import (
	"context"
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

	// Логируем начало
	log.Println("📨 === REGISTER REQUEST START ===")

	// Разбираем запрос
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ JSON decode error: %v", err)
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Логируем данные
	log.Printf("📝 Регистрация: email=%s, username=%s, full_name='%s'",
		req.Email, req.Username, req.FullName)

	// Валидируем данные
	if err := validateRegisterRequest(req); err != nil {
		log.Printf("❌ Validation error: %v", err)
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)"

	log.Printf("🔍 Проверка существования: email=%s, username=%s",
		strings.ToLower(req.Email), req.Username)

	err := h.DB.QueryRowContext(ctx, query, strings.ToLower(req.Email), req.Username).Scan(&exists)
	if err != nil {
		log.Printf("❌ Database error checking user existence: %v", err)
		sendJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("🔍 Результат проверки: exists=%v", exists)

	if exists {
		log.Printf("⚠️ Пользователь уже существует! Возвращаем 409")

		// Найдем конфликтующего пользователя для отладки
		var conflictEmail, conflictUsername string
		h.DB.QueryRowContext(ctx,
			"SELECT email, username FROM users WHERE email = $1 OR username = $2 LIMIT 1",
			strings.ToLower(req.Email), req.Username,
		).Scan(&conflictEmail, &conflictUsername)

		log.Printf("⚠️ Конфликт с: email='%s', username='%s'", conflictEmail, conflictUsername)

		sendJSONError(w, "User with this email or username already exists", http.StatusConflict)
		return
	}

	log.Printf("✅ Пользователь не существует, создаем...")

	// Хешируем пароль
	passwordHash, err := models.HashPassword(req.Password)
	if err != nil {
		log.Printf("❌ Password hashing error: %v", err)
		sendJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базе данных
	var id string
	var createdAt, updatedAt time.Time

	// Явно генерируем UUID и обрабатываем NULL для full_name
	insertQuery := `
		INSERT INTO users (
			id,
			email, 
			username, 
			password_hash, 
			full_name, 
			score, 
			current_streak, 
			best_streak
		) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, created_at, updated_at
	`

	var fullNameParam *string
	if req.FullName != nil {
		fullNameParam = req.FullName
		log.Printf("full_name будет: '%s'", *req.FullName)
	} else {
		fullNameParam = nil
		log.Printf("full_name будет: NULL")
	}

	log.Printf("📤 Выполняем INSERT с параметрами: email=%s, username=%s, full_name=%v",
		strings.ToLower(req.Email), req.Username, fullNameParam)

	err = h.DB.QueryRowContext(ctx, insertQuery,
		strings.ToLower(req.Email),
		req.Username,
		passwordHash,
		fullNameParam,
		0,
		0,
		0,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		log.Printf("❌ Database insert error: %v", err)
		log.Printf("❌ Детали ошибки: email=%s, username=%s", req.Email, req.Username)

		// Проверим структуру таблицы
		log.Println("🔍 Проверка структуры таблицы users...")
		h.debugTableStructure(ctx)

		sendJSONError(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Пользователь создан! ID: %s, email: %s", id, req.Email)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	log.Println("📨 === REGISTER REQUEST SUCCESS ===")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Error encoding response: %v", err)
	}
}

// Вспомогательная функция для отладки структуры таблицы
func (h *AuthHandler) debugTableStructure(ctx context.Context) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns 
		WHERE table_name = 'users' 
		ORDER BY ordinal_position
	`)

	if err != nil {
		log.Printf("❌ Не удалось получить структуру таблицы: %v", err)
		return
	}
	defer rows.Close()

	log.Println("📊 Структура таблицы 'users':")
	log.Println("   Колонка           | Тип        | Nullable | Default")
	log.Println("   ------------------|------------|----------|---------")

	for rows.Next() {
		var colName, dataType, isNullable, columnDefault sql.NullString
		rows.Scan(&colName, &dataType, &isNullable, &columnDefault)

		def := "-"
		if columnDefault.Valid {
			def = columnDefault.String
		}

		log.Printf("   %-18s | %-10s | %-8s | %s",
			colName.String, dataType.String, isNullable.String, def)
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
	log.Printf("❌ Отправляем ошибку: %s (код: %d)", message, statusCode)

	response := models.APIResponse{
		Status: "error",
		Error:  message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Error encoding error response: %v", err)
	}
}
