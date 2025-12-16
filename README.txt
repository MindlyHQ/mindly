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