package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mindly/api/internal/database"  
    "github.com/mindly/api/internal/handlers"
)

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого origin (для разработки)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Обрабатываем предварительные OPTIONS запросы
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос дальше
		next.ServeHTTP(w, r)
	})
}

// JSON middleware
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("🚀 Starting Mindly API Server...")

	// Создаем контекст
	ctx := context.Background()
	
	// Подключаемся к базе данных
	cfg := database.DefaultConfig()
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Database connected successfully")

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(db)

	// Настраиваем маршруты
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("GET /health", healthHandler)
	
	// Auth endpoints
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	
	// Добавляем middleware
	handler := enableCORS(mux)

	// Настраиваем сервер
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("🌐 Server listening on http://localhost%s", server.Addr)
		log.Printf("📊 Health check: http://localhost%s/health", server.Addr)
		log.Printf("👤 Register endpoint: POST http://localhost%s/api/auth/register", server.Addr)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	// Блокируемся до получения сигнала
	sig := <-stop
	log.Printf("🛑 Received signal: %v", sig)
	log.Println("Shutting down server...")
	
	// Создаем контекст с таймаутом для graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Server shutdown error: %v", err)
	}
	
	log.Println("👋 Server stopped gracefully")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":   "ok",
		"service":  "mindly-api",
		"version":  "1.0.0",
		"database": "connected",
		"time":     time.Now().UTC().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
}
