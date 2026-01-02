package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🌱 Загрузка тестовых данных для Mindly LearnStream...")

	// Подключение к БД (используем твои настройки)
	connStr := "host=localhost port=5432 user=mindly password=mindly123 dbname=mindly_dev sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ БД не отвечает: %v", err)
	}

	fmt.Println("✅ Подключение к БД успешно")

	// 1. ПРОВЕРЯЕМ И СОЗДАЕМ ПОЛЬЗОВАТЕЛЯ (если нет)
	fmt.Println("\n1. Работа с пользователем...")

	var userID int

	// Проверяем, есть ли пользователи
	err = db.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаём нового пользователя (только email и created_at, как в твоей схеме)
			fmt.Println("   👤 Создаём нового пользователя...")
			err = db.QueryRowContext(ctx,
				"INSERT INTO users (email) VALUES ($1) RETURNING id",
				"demo@mindly.ru",
			).Scan(&userID)

			if err != nil {
				log.Printf("❌ Не удалось создать пользователя: %v", err)
				// Пробуем создать с другим email
				db.QueryRowContext(ctx,
					"INSERT INTO users (email) VALUES ($1) RETURNING id",
					"test@mindly.ru",
				).Scan(&userID)
			}
		} else {
			log.Printf("⚠️ Ошибка при проверке пользователей: %v", err)
			return
		}
	}

	fmt.Printf("   👤 Используем User ID: %d\n", userID)

	// 2. СОЗДАЕМ АВТОРА (эксперта)
	fmt.Println("\n2. Создаём автора-эксперта...")

	var authorID string

	// Проверяем, есть ли уже авторы
	err = db.QueryRowContext(ctx, "SELECT id FROM authors LIMIT 1").Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаём нового автора
			authorQuery := `
                INSERT INTO authors (user_id, full_name, expertise_area, trust_tier) 
                VALUES ($1, $2, $3, $4)
                RETURNING id
            `

			err = db.QueryRowContext(ctx, authorQuery,
				userID,
				"Дмитрий Программист",
				"IT",
				"gold",
			).Scan(&authorID)

			if err != nil {
				log.Printf("❌ Не удалось создать автора: %v", err)
				return
			}
		} else {
			log.Printf("⚠️ Ошибка при проверке авторов: %v", err)
			return
		}
	}

	fmt.Printf("   📝 Используем Author ID: %s\n", authorID)

	// 3. ДОБАВЛЯЕМ ТЕСТОВЫЕ ВИДЕО
	fmt.Println("\n3. Добавляем тестовые видео...")

	// Сначала удалим старые тестовые видео (если есть)
	_, _ = db.ExecContext(ctx, "DELETE FROM videos WHERE title LIKE 'Тест:%' OR title LIKE '%API%'")

	// Тестовые видео с публичными ссылками
	testVideos := []struct {
		title, description, videoURL, thumbnailURL string
		durationSec                                int
		tags                                       []string
	}{
		{
			"Что такое API за 60 секунд",
			"Простое объяснение API для начинающих. Как приложения общаются между собой.",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
			"https://img.youtube.com/vi/s7wmiS2mSXY/mqdefault.jpg",
			60,
			[]string{"программирование", "api", "веб"},
		},
		{
			"Основы Go языка",
			"Почему Go такой популярный для backend-разработки. Основные фичи.",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
			"https://img.youtube.com/vi/yoTahYcKnyo/mqdefault.jpg",
			90,
			[]string{"golang", "go", "программирование"},
		},
		{
			"HTTP и HTTPS простыми словами",
			"В чём разница между HTTP и HTTPS. Почему важно использовать HTTPS.",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
			"https://img.youtube.com/vi/hExRDVZHhig/mqdefault.jpg",
			75,
			[]string{"http", "безопасность", "веб"},
		},
		{
			"Базы данных: SQL за 80 секунд",
			"Основные понятия SQL для начинающих. SELECT, INSERT, UPDATE, DELETE.",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
			"https://img.youtube.com/vi/7V_mN1-d2eM/mqdefault.jpg",
			80,
			[]string{"sql", "базы данных", "postgresql"},
		},
		{
			"Как победить прокрастинацию",
			"Практические советы для разработчиков и не только.",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
			"https://img.youtube.com/vi/Qvcx7Y4caQE/mqdefault.jpg",
			85,
			[]string{"продуктивность", "психология", "саморазвитие"},
		},
	}

	videosAdded := 0
	videoIDs := []string{}

	for i, video := range testVideos {
		var videoID string
		videoQuery := `
            INSERT INTO videos (author_id, title, description, video_url, thumbnail_url, duration_sec, tags) 
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id
        `

		// Преобразуем массив тегов в формат PostgreSQL
		tagsStr := "{" + strings.Join(video.tags, ",") + "}"

		err = db.QueryRowContext(ctx, videoQuery,
			authorID,
			video.title,
			video.description,
			video.videoURL,
			video.thumbnailURL,
			video.durationSec,
			tagsStr,
		).Scan(&videoID)

		if err != nil {
			log.Printf("⚠️ Ошибка при добавлении видео '%s': %v", video.title, err)
			// Пробуем без thumbnail_url
			db.QueryRowContext(ctx, `
                INSERT INTO videos (author_id, title, description, video_url, duration_sec, tags) 
                VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id
            `,
				authorID,
				video.title,
				video.description,
				video.videoURL,
				video.durationSec,
				tagsStr,
			).Scan(&videoID)
		}

		if videoID != "" {
			videosAdded++
			videoIDs = append(videoIDs, videoID)
			fmt.Printf("   ✅ Видео %d: %s (ID: %s)\n", i+1, video.title, videoID[:8])
		}
	}

	// 4. ДОБАВЛЯЕМ ТЕСТЫ К ВИДЕО
	fmt.Println("\n4. Добавляем тесты к видео...")

	testsAdded := 0
	for _, videoID := range videoIDs {
		quizQuery := `
            INSERT INTO quizzes (video_id, question, correct_answer, wrong_answers) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (video_id) DO NOTHING
        `

		_, err := db.ExecContext(ctx, quizQuery,
			videoID,
			"Был ли этот материал полезен?",
			"Да, узнал что-то новое",
			`{"Уже знал это", "Слишком сложно", "Не по теме"}`,
		)

		if err == nil {
			testsAdded++
		}
	}

	// 5. ФИНАЛЬНАЯ ПРОВЕРКА
	fmt.Println("\n5. Финальная проверка...")

	var videoCount int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM videos").Scan(&videoCount)

	var authorCount int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM authors").Scan(&authorCount)

	fmt.Printf("   📊 В БД теперь: %d видео, %d авторов\n", videoCount, authorCount)

	if videosAdded > 0 {
		fmt.Printf("\n🎉 УСПЕХ! Загружено %d видео и %d тестов.\n", videosAdded, testsAdded)
		fmt.Println("🔗 Проверьте API: http://localhost:8081/api/feed?limit=5")
		fmt.Println("📺 Пример видео URL: https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4")
	} else {
		fmt.Println("\n⚠️ Видео не были добавлены. Проверьте структуру таблиц.")
		fmt.Println("   Выполните: \\dt в psql для проверки таблиц")
	}
}
