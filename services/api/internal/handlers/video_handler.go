package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/mindly/api/internal/database"
)

type VideoHandler struct {
    videoRepo *database.VideoRepository
}

// Принимаем *sql.DB напрямую
func NewVideoHandler(db *sql.DB) *VideoHandler {
    videoRepo := database.NewVideoRepository(db)
    return &VideoHandler{videoRepo: videoRepo}
}

func (h *VideoHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
    // Получаем user_id из query-параметра
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        userID = "temp_user_id"
    }
    
    // Получаем лимит
    limitStr := r.URL.Query().Get("limit")
    if limitStr == "" {
        limitStr = "10"
    }
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }
    if limit > 50 {
        limit = 50
    }
    
    // Получаем видео из репозитория
    videos, err := h.videoRepo.GetFeed(r.Context(), userID, limit)
    if err != nil {
        http.Error(w, `{"error": "Не удалось загрузить ленту: `+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    
    // Формируем ответ
    response := map[string]interface{}{
        "success": true,
        "data":    videos,
        "count":   len(videos),
    }
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, `{"error": "Ошибка кодирования ответа"}`, http.StatusInternalServerError)
        return
    }
}