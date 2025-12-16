package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
    app := fiber.New()
    
    app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000,http://localhost:8081",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin,Content-Type,Accept,Authorization",
    }))
    
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":  "ok",
            "service": "mindly-api",
            "version": "0.1.0",
        })
    })
    
    app.Get("/api/test", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Mindly API is working!",
        })
    })
    
    log.Println("🚀 Mindly API starting on http://localhost:8080")
    if err := app.Listen(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}