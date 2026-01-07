@echo off
chcp 65001 >nul
cls
echo Очистка и запуск Mindly...
echo.

:: 1. ОСТАНОВИТЬ ВСЁ
docker-compose down 2>nul
docker stop amwesome_goldwasser 2>nul
docker rm amwesome_goldwasser 2>nul

:: 2. ЗАПУСТИТЬ FRESH
docker-compose up -d
timeout /t 5 /nobreak >nul

:: 3. ЗАПУСТИТЬ БЭКЕНД
start /B cmd /c "cd services/api && go run ./cmd/api/main.go"
timeout /t 3 /nobreak >nul
cd services/api && go run ./scripts/seed.go && cd ../..

:: 4. ЗАПУСТИТЬ ПРИЛОЖЕНИЕ С QR
cd apps/mobile
cls
echo ========================
echo    QR-КОД ДЛЯ EXPO GO
echo ========================
echo.
npm start