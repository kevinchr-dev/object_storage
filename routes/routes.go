package routes

import (
	"object-storage-server/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, fileHandler *handlers.FileHandler) {
	// API routes
	api := app.Group("/api")

	// File operations
	api.Post("/upload", fileHandler.UploadFile)
	api.Get("/files/:filename", fileHandler.DownloadFile)
	api.Get("/files/view/:filename", fileHandler.ViewFile)
	api.Get("/files/info/:filename", fileHandler.GetFileInfo)         // Deprecated
	api.Get("/files/metadata/:filename", fileHandler.GetFileMetadata) // New metadata endpoint

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Object Storage Server is running",
		})
	})
}
