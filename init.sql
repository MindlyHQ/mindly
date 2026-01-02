-- 1. ПОЛЬЗОВАТЕЛИ (основа системы)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    score INTEGER DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    best_streak INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. АВТОРЫ (эксперты, создатели контента)
CREATE TABLE authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Связь с учетной записью пользователя (автор тоже может смотреть видео)
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    expertise_area VARCHAR(100), -- 'IT', 'Психология', 'Финансы'
    -- Уровень доверия: для презентации можно назначать вручную или по упрощенной логике
    trust_tier VARCHAR(20) DEFAULT 'silver',
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. ВИДЕО (главная сущность)
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    -- Ссылка на облачное хранилище (например, Vimeo, YouTube, S3)
    video_url TEXT NOT NULL,
    thumbnail_url TEXT,
    -- Ограничение по длине: 20-90 секунд, как в плане
    duration_sec INTEGER NOT NULL CHECK (duration_sec BETWEEN 20 AND 300),
    -- Теги для рекомендаций (например, ['python', 'программирование', 'для начинающих'])
    tags TEXT[] DEFAULT '{}',
    -- Для презентации: статус модерации можно ставить 'approved' вручную
    moderation_status VARCHAR(20) DEFAULT 'approved',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. ВОПРОСЫ (интерактив после видео)
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Один вопрос к одному видео
    video_id UUID UNIQUE NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    -- Для простоты храним один правильный и несколько неправильных вариантов
    correct_answer VARCHAR(255) NOT NULL,
    wrong_answers VARCHAR(255)[] DEFAULT '{}',
    -- Баллы за правильный ответ (основа геймификации)
    points_awarded INTEGER DEFAULT 1
);

-- 5. ПРОГРЕСС ПОЛЬЗОВАТЕЛЯ (самая важная таблица для аналитики и геймификации)
CREATE TABLE user_video_progress (
    -- Составной первичный ключ: один пользователь - одно видео
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    -- Факт просмотра
    is_watched BOOLEAN DEFAULT FALSE,
    watched_at TIMESTAMP,
    -- Результат прохождения теста
    quiz_attempted BOOLEAN DEFAULT FALSE,
    quiz_correct BOOLEAN DEFAULT FALSE,
    -- Начисленные баллы за это видео
    points_earned INTEGER DEFAULT 0,
    -- КРИТИЧЕСКИ ВАЖНОЕ поле для подсчета серий (streak)
    -- Будем считать streak по дням, когда было любое взаимодействие
    interaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (user_id, video_id)
);

-- 6. ОБЩИЕ БАЛЛЫ И СТАТИСТИКА ПОЛЬЗОВАТЕЛЯ (для личного кабинета)
CREATE TABLE user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    -- Суммарные баллы (можно пересчитывать из user_video_progress, но для скорости храним отдельно)
    total_points INTEGER DEFAULT 0,
    -- Текущая серия дней (streak)
    current_streak_days INTEGER DEFAULT 0,
    -- Максимальная достигнутая серия
    max_streak_days INTEGER DEFAULT 0,
    -- Дата последней активности для проверки сброса streak
    last_activity_date DATE DEFAULT CURRENT_DATE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения ключевых запросов (лента, прогресс)
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX idx_videos_tags ON videos USING GIN(tags);
CREATE INDEX idx_user_progress_user_interaction ON user_video_progress(user_id, interaction_date DESC);
CREATE INDEX idx_user_progress_video ON user_video_progress(video_id);