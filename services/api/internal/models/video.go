package models

import "time"

type Video struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    VideoURL    string    `json:"video_url"`
    ThumbnailURL string   `json:"thumbnail_url,omitempty"`
    DurationSec int       `json:"duration_sec"`
    Tags        []string  `json:"tags"`
    CreatedAt   time.Time `json:"created_at"`
}

type Author struct {
    ID           string `json:"id"`
    FullName     string `json:"full_name"`
    ExpertiseArea string `json:"expertise_area"`
    TrustTier    string `json:"trust_tier"`
    IsVerified   bool   `json:"is_verified"`
}

// VideoWithAuthor - объединенные данные для ленты
type VideoWithAuthor struct {
    Video  `json:",inline"`
    Author Author `json:"author"`
}

type Quiz struct {
    Question      string   `json:"question"`
    CorrectAnswer string   `json:"correct_answer"`
    WrongAnswers  []string `json:"wrong_answers"`
    PointsAwarded int      `json:"points_awarded"`
}