package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"object-storage-server/config"
	_ "object-storage-server/docs" // Swagger docs
	"object-storage-server/handlers"
	"object-storage-server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default configuration")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Ensure upload directory exists
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	// Create Fiber app with optimized settings for concurrent connections
	app := fiber.New(fiber.Config{
		BodyLimit:             int(cfg.MaxFileSize), // 4GB max per request
		Concurrency:           256 * 1024,           // Handle up to 256k concurrent connections
		ReadBufferSize:        16384,                // 16KB read buffer for large files
		WriteBufferSize:       16384,                // 16KB write buffer for large files
		StreamRequestBody:     true,                 // Stream large file uploads
		DisableStartupMessage: false,
		Prefork:               false, // Set to true for multi-process mode (use with caution)
		EnablePrintRoutes:     false,
		ReadTimeout:           0, // No timeout for large file uploads
		WriteTimeout:          0, // No timeout for large file downloads
	})

	// Middleware
	// 1. Recovery from panics
	app.Use(recover.New())

	// 2. Logger
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency}) [${ip}]\n",
	}))

	// 3. Compression for responses (gzip)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Balance between speed and compression
	}))

	// 4. Rate limiter to prevent abuse (100 requests per minute per IP)
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit per IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too many requests, please try again later",
			})
		},
	}))

	// 5. CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedHosts,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Initialize handlers
	fileHandler := handlers.NewFileHandler(cfg)

	// Setup routes
	routes.SetupRoutes(app, fileHandler)

	// Swagger documentation - must be after routes
	app.Get("/docs/*", swagger.New(swagger.Config{
		DeepLinking:  true,
		DocExpansion: "list",
		Title:        "Object Storage API Documentation",
	}))

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Object Storage Server running on %s", cfg.BaseURL)
	log.Printf("📁 Upload directory: %s", cfg.UploadDir)
	log.Printf("📊 Max file size: %d bytes (%.2f MB)", cfg.MaxFileSize, float64(cfg.MaxFileSize)/(1024*1024))

	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
