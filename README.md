# Mindly

Educational platform in TikTok format.

## What is Mindly?
Short educational videos + instant knowledge check.

## Tech Stack
- React Native (mobile app)
- Go (backend server)
- PostgreSQL (database)
- Redis (cache)

## Quick Start
```bash
# Start infrastructure
docker-compose up -d

# Run backend
cd services/api && go run cmd/api/main.go

# Run mobile app
cd apps/mobile && npx expo start



main test 
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=mindly password=mindly123 dbname=mindly_dev sslmode=disable"

	log.Println("=== ТЕСТ ПОРТ 5432 ===")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ OPEN:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("❌ PING:", err)
	}

	log.Println("✅ УСПЕХ!")
}

ИНСТРУКЦИЯ К ЗАПУСКУ
в трех разных терминалах: 
1:C:\Users\denis\Documents\Projects\mindly\services\api>go run ./cmd/api/main.go
2:PS C:\Users\denis\Documents\Projects\mindly\apps\mobile> npm start
3:PS C:\Users\denis\Documents\Projects\mindly\services\api> go run ./scripts/seed.go

PostgreSQL: 5432

Redis: 6379

API сервер: 8081

Metro Bundler: 8082