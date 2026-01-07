package main

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func safeShortID(id string, length int) string {
	if id == "" {
		return "НЕТ ID"
	}
	if len(id) >= length {
		return id[:length] + "..."
	}
	return id
}

func generatePasswordHash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

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

	var userID string

	// Проверяем существующих пользователей
	err = db.QueryRowContext(ctx, "SELECT id::text FROM users LIMIT 1").Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаём нового пользователя с ВСЕМИ обязательными полями
			fmt.Println("   👤 Создаём нового пользователя...")

			// Генерируем хеш пароля
			passwordHash := generatePasswordHash("mindly123")
			currentTime := time.Now()

			// Пытаемся создать пользователя со всеми обязательными полями
			err = db.QueryRowContext(ctx,
				`INSERT INTO users (
					email, 
					username, 
					password_hash, 
					full_name,
					score,
					current_streak,
					best_streak,
					created_at,
					updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
				RETURNING id::text`,
				"demo@mindly.ru",
				"demo_user",
				passwordHash,
				"Демо Пользователь",
				0, // score
				0, // current_streak
				0, // best_streak
				currentTime,
				currentTime,
			).Scan(&userID)

			if err != nil {
				log.Printf("   ⚠️ Не удалось создать первого пользователя: %v", err)

				// Пробуем минимальный набор полей
				err = db.QueryRowContext(ctx,
					`INSERT INTO users (
						email, 
						username, 
						password_hash,
						created_at,
						updated_at
					) VALUES ($1, $2, $3, $4, $5) 
					RETURNING id::text`,
					"test@mindly.ru",
					"test_user",
					generatePasswordHash("test123"),
					currentTime,
					currentTime,
				).Scan(&userID)

				if err != nil {
					log.Printf("❌ Не удалось создать пользователя: %v", err)

					// Ещё одна попытка с другим именем пользователя
					err = db.QueryRowContext(ctx,
						`INSERT INTO users (
							email, 
							username, 
							password_hash,
							created_at,
							updated_at
						) VALUES ($1, $2, $3, $4, $5) 
						RETURNING id::text`,
						"admin@mindly.ru",
						"admin",
						generatePasswordHash("admin123"),
						currentTime,
						currentTime,
					).Scan(&userID)

					if err != nil {
						log.Fatalf("❌ Все попытки создать пользователя не удались: %v", err)
					}
				}
			}
		} else {
			log.Printf("⚠️ Ошибка при проверке пользователей: %v", err)
			return
		}
	}

	fmt.Printf("   👤 Используем User ID: %s\n", safeShortID(userID, 8))

	// 2. СОЗДАЕМ АВТОРА (эксперта)
	fmt.Println("\n2. Создаём автора-эксперта...")

	var authorID string

	// Проверяем, есть ли уже авторы
	err = db.QueryRowContext(ctx, "SELECT id::text FROM authors LIMIT 1").Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаём нового автора
			authorQuery := `
                INSERT INTO authors (
					user_id, 
					full_name, 
					expertise_area, 
					trust_tier,
					bio,
					created_at,
					updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id::text
            `

			currentTime := time.Now()
			err = db.QueryRowContext(ctx, authorQuery,
				userID,
				"Дмитрий Программист",
				"IT",
				"gold",
				"Опытный разработчик с 10-летним стажем. Специализируется на Go, микросервисах и DevOps.",
				currentTime,
				currentTime,
			).Scan(&authorID)

			if err != nil {
				log.Printf("❌ Не удалось создать автора: %v", err)

				// Пробуем без optional полей
				db.QueryRowContext(ctx, `
                    INSERT INTO authors (
						user_id, 
						full_name, 
						expertise_area
					) VALUES ($1, $2, $3)
                    RETURNING id::text
                `,
					userID,
					"Дмитрий Программист",
					"IT",
				).Scan(&authorID)

				if err != nil {
					log.Printf("❌ Вторая попытка создать автора тоже не удалась: %v", err)

					// Посмотрим структуру таблицы authors
					var authorColumns string
					db.QueryRowContext(ctx,
						`SELECT string_agg(column_name || ' ' || 
							CASE WHEN is_nullable = 'NO' THEN 'NOT NULL' ELSE '' END, ', ') 
						 FROM information_schema.columns WHERE table_name = 'authors'`).Scan(&authorColumns)
					if authorColumns != "" {
						log.Printf("   Структура таблицы authors: %s\n", authorColumns)
					}
					return
				}
			}
		} else {
			log.Printf("⚠️ Ошибка при проверке авторов: %v", err)
			return
		}
	}

	fmt.Printf("   📝 Используем Author ID: %s\n", safeShortID(authorID, 8))

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
	currentTime := time.Now()

	for i, video := range testVideos {
		var videoID string
		videoQuery := `
            INSERT INTO videos (
				author_id, 
				title, 
				description, 
				video_url, 
				thumbnail_url, 
				duration_sec, 
				tags,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id::text
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
			currentTime,
			currentTime,
		).Scan(&videoID)

		if err != nil {
			log.Printf("⚠️ Ошибка при добавлении видео '%s': %v", video.title, err)

			// Пробуем без thumbnail_url
			err = db.QueryRowContext(ctx, `
                INSERT INTO videos (
					author_id, 
					title, 
					description, 
					video_url, 
					duration_sec, 
					tags,
					created_at,
					updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                RETURNING id::text
            `,
				authorID,
				video.title,
				video.description,
				video.videoURL,
				video.durationSec,
				tagsStr,
				currentTime,
				currentTime,
			).Scan(&videoID)

			if err != nil {
				log.Printf("❌ Вторая попытка для видео '%s' тоже не удалась: %v", video.title, err)

				// Посмотрим структуру таблицы videos
				var videoColumns string
				db.QueryRowContext(ctx,
					`SELECT string_agg(column_name || ' ' || 
						CASE WHEN is_nullable = 'NO' THEN 'NOT NULL' ELSE '' END, ', ') 
					 FROM information_schema.columns WHERE table_name = 'videos'`).Scan(&videoColumns)
				if videoColumns != "" {
					log.Printf("   Структура таблицы videos: %s\n", videoColumns)
				}
				continue
			}
		}

		if videoID != "" {
			videosAdded++
			videoIDs = append(videoIDs, videoID)
			fmt.Printf("   ✅ Видео %d: %s (ID: %s)\n", i+1, video.title, safeShortID(videoID, 8))
		}
	}

	// 4. ДОБАВЛЯЕМ ТЕСТЫ К ВИДЕО (если есть таблица quizzes)
	fmt.Println("\n4. Проверяем таблицу quizzes...")

	var tableExists bool
	db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'quizzes')").Scan(&tableExists)

	testsAdded := 0
	if tableExists && len(videoIDs) > 0 {
		fmt.Println("   Добавляем тесты к видео...")

		for _, videoID := range videoIDs {
			if videoID == "" {
				continue
			}

			quizQuery := `
				INSERT INTO quizzes (
					video_id, 
					question, 
					correct_answer, 
					wrong_answers,
					created_at,
					updated_at
				) VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (video_id) DO NOTHING
			`

			_, err := db.ExecContext(ctx, quizQuery,
				videoID,
				"Был ли этот материал полезен?",
				"Да, узнал что-то новое",
				`{"Уже знал это", "Слишком сложно", "Не по теме"}`,
				currentTime,
				currentTime,
			)

			if err != nil {
				log.Printf("⚠️ Ошибка при добавлении теста для видео %s: %v", safeShortID(videoID, 8), err)
			} else {
				testsAdded++
			}
		}
	} else {
		fmt.Println("   ⚠️ Таблица quizzes не существует или нет видео")
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

		// Дополнительная информация о пользователе
		var userEmail, userName string
		db.QueryRowContext(ctx, "SELECT email, username FROM users WHERE id = $1", userID).Scan(&userEmail, &userName)
		fmt.Printf("👤 Тестовый пользователь: %s (%s)\n", userName, userEmail)
		fmt.Println("🔐 Пароль: mindly123 (или test123/admin123 в зависимости от созданного)")
	} else {
		fmt.Println("\n⚠️ Видео не были добавлены. Возможные причины:")
		fmt.Println("   • Проверьте структуру таблицы videos")
		fmt.Println("   • Проверьте подключение к БД")
		fmt.Println("   • Убедитесь, что author_id корректен")

		// Проверим, есть ли таблица videos
		var videosTableExists bool
		db.QueryRowContext(ctx,
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'videos')").Scan(&videosTableExists)

		if !videosTableExists {
			fmt.Println("   ❌ Таблица videos не существует!")
		} else {
			fmt.Println("   ✅ Таблица videos существует")
		}
	}
}
