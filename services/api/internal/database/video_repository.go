package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mindly/api/internal/models"
)

type VideoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) GetFeed(ctx context.Context, userID string, limit int) ([]models.VideoWithAuthor, error) {
    query := `
        SELECT 
            v.id, v.title, v.description, v.video_url, v.thumbnail_url,
            v.duration_sec, v.tags, v.created_at,
            a.id, a.full_name, a.expertise_area, a.trust_tier, a.is_verified
        FROM videos v
        JOIN authors a ON v.author_id = a.id
        WHERE v.moderation_status = 'approved'
        ORDER BY v.created_at DESC
        LIMIT $1
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit)
    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    defer rows.Close()
    
    var videos []models.VideoWithAuthor
    
    for rows.Next() {
        var v models.VideoWithAuthor
        var tagsRaw []byte // Изменяем тип на []byte для обработки массива
        var tagsStr string
        var thumbnailURL sql.NullString // Обработка NULL значения
        
        err := rows.Scan(
            &v.ID, &v.Title, &v.Description, &v.VideoURL, &thumbnailURL,
            &v.DurationSec, &tagsRaw, &v.CreatedAt,
            &v.Author.ID, &v.Author.FullName, &v.Author.ExpertiseArea,
            &v.Author.TrustTier, &v.Author.IsVerified,
        )
        if err != nil {
            return nil, fmt.Errorf("scan error: %w", err)
        }
        
        // Обработка thumbnail_url (может быть NULL)
        if thumbnailURL.Valid {
            v.ThumbnailURL = thumbnailURL.String
        }
        
        // Преобразуем массив tags из PostgreSQL формата в []string
        tagsStr = string(tagsRaw)
        v.Tags = parsePostgresArray(tagsStr)
        v.Video.Tags = v.Tags // Копируем также в embedded структуру
        
        videos = append(videos, v)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }
    
    return videos, nil
}

// Вспомогательная функция для парсинга PostgreSQL массива
func parsePostgresArray(arrayStr string) []string {
    // Убираем фигурные скобки
    arrayStr = strings.Trim(arrayStr, "{}")
    
    if arrayStr == "" {
        return []string{}
    }
    
    // Разделяем по запятым, учитывая кавычки
    var result []string
    var current strings.Builder
    inQuotes := false
    
    for i := 0; i < len(arrayStr); i++ {
        ch := arrayStr[i]
        
        switch {
        case ch == '"':
            if i+1 < len(arrayStr) && arrayStr[i+1] == '"' {
                // Экранированная кавычка ""
                current.WriteByte('"')
                i++ // Пропускаем следующую кавычку
            } else {
                inQuotes = !inQuotes
            }
        case ch == ',' && !inQuotes:
            result = append(result, current.String())
            current.Reset()
        default:
            current.WriteByte(ch)
        }
    }
    
    // Добавляем последний элемент
    if current.Len() > 0 {
        result = append(result, current.String())
    }
    
    return result
}
