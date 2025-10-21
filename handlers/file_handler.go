package handlers

import (
	"fmt"
	"io"
	"object-storage-server/config"
	"object-storage-server/models"
	"object-storage-server/utils"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	Config *config.Config
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{Config: cfg}
}

// UploadFile handles file upload
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Message: "No file uploaded or invalid form data",
		})
	}

	// Check file size
	if file.Size > h.Config.MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Message: fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", h.Config.MaxFileSize),
		})
	}

	// Generate unique filename with UUID v7
	uniqueFileName := utils.GenerateUniqueFileName(file.Filename)

	// Create full path
	fullPath := filepath.Join(h.Config.UploadDir, uniqueFileName)

	// Save file
	if err := c.SaveFile(file, fullPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Message: "Failed to save file",
		})
	}

	// Check file type
	isImage := utils.IsImage(uniqueFileName)
	isVideo := utils.IsVideo(uniqueFileName)
	isAudio := utils.IsAudio(uniqueFileName)

	fileType := "other"
	if isImage {
		fileType = "image"
	} else if isVideo {
		fileType = "video"
	} else if isAudio {
		fileType = "audio"
	}

	// Prepare response
	response := models.UploadResponse{
		Success:     true,
		Message:     "File uploaded successfully",
		FileName:    uniqueFileName,
		FileURL:     fmt.Sprintf("%s/api/files/%s", h.Config.BaseURL, uniqueFileName),
		MetadataURL: fmt.Sprintf("%s/api/files/metadata/%s", h.Config.BaseURL, uniqueFileName),
		FileSize:    file.Size,
		FileType:    fileType,
		IsImage:     isImage,
		IsVideo:     isVideo,
		IsAudio:     isAudio,
	}

	// Generate view URLs
	viewURLs := &models.ViewURLs{
		Original: fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, uniqueFileName),
	}

	// Get global worker pool for background processing
	workerPool := utils.GetWorkerPool()

	// If image, create resized versions (non-blocking for large images)
	if isImage {
		// For small images (< 2MB), process synchronously for instant response
		if file.Size < 2*1024*1024 {
			resizedFiles, err := utils.ResizeImage(fullPath, h.Config.UploadDir, uniqueFileName)
			if err == nil && len(resizedFiles) > 0 {
				// Add resized version URLs
				if thumbnail, ok := resizedFiles["thumbnail"]; ok {
					viewURLs.Thumbnail = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, thumbnail)
				}
				if small, ok := resizedFiles["small"]; ok {
					viewURLs.Small = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, small)
				}
				if medium, ok := resizedFiles["medium"]; ok {
					viewURLs.Medium = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, medium)
				}
				if large, ok := resizedFiles["large"]; ok {
					viewURLs.Large = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, large)
				}
			}
		} else {
			// For large images (>= 2MB), process in worker pool
			workerPool.Submit(utils.Job{
				Type:      "image",
				FilePath:  fullPath,
				UploadDir: h.Config.UploadDir,
				FileName:  uniqueFileName,
			})
			response.Message = "File uploaded successfully. Image processing in progress..."
		}
	}

	// If video, create multiple resolutions and thumbnail
	if isVideo && utils.CheckFFmpegInstalled() {
		// Submit to worker pool for controlled concurrent processing
		workerPool.Submit(utils.Job{
			Type:      "video",
			FilePath:  fullPath,
			UploadDir: h.Config.UploadDir,
			FileName:  uniqueFileName,
		})
		// Note: Video processing happens in worker pool
		response.Message = "File uploaded successfully. Video processing queued..."
	}

	// If audio, create multiple bitrates
	if isAudio && utils.CheckFFmpegInstalled() {
		// Submit to worker pool for controlled concurrent processing
		workerPool.Submit(utils.Job{
			Type:      "audio",
			FilePath:  fullPath,
			UploadDir: h.Config.UploadDir,
			FileName:  uniqueFileName,
		})
		// Note: Audio processing happens in worker pool
		response.Message = "File uploaded successfully. Audio processing queued..."
	}

	response.ViewURLs = viewURLs

	return c.Status(fiber.StatusOK).JSON(response)
}

// DownloadFile handles file download
func (h *FileHandler) DownloadFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Message: "Filename is required",
		})
	}

	// Prevent directory traversal
	filename = filepath.Base(filename)

	// Create full path
	fullPath := filepath.Join(h.Config.UploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
	}

	// Set content disposition header for download
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Send file
	return c.SendFile(fullPath)
}

// ViewFile handles file viewing (inline)
func (h *FileHandler) ViewFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Message: "Filename is required",
		})
	}

	// Prevent directory traversal
	filename = filepath.Base(filename)

	// Create full path
	fullPath := filepath.Join(h.Config.UploadDir, filename)

	// Check if file exists
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Message: "Failed to open file",
		})
	}
	defer file.Close()

	// Set content type
	contentType := utils.GetContentType(filename)
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	c.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Stream file
	_, err = io.Copy(c.Response().BodyWriter(), file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Message: "Failed to stream file",
		})
	}

	return nil
}

// GetFileInfo returns file information (deprecated, use GetFileMetadata)
func (h *FileHandler) GetFileInfo(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Message: "Filename is required",
		})
	}

	// Prevent directory traversal
	filename = filepath.Base(filename)

	// Create full path
	fullPath := filepath.Join(h.Config.UploadDir, filename)

	// Check if file exists
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":   true,
		"file_name": filename,
		"file_size": fileInfo.Size(),
		"modified":  fileInfo.ModTime(),
	})
}

// GetFileMetadata returns detailed file metadata including all URLs
func (h *FileHandler) GetFileMetadata(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Message: "Filename is required",
		})
	}

	// Prevent directory traversal
	filename = filepath.Base(filename)

	// Create full path
	fullPath := filepath.Join(h.Config.UploadDir, filename)

	// Check if file exists
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
	}

	isImage := utils.IsImage(filename)
	isVideo := utils.IsVideo(filename)
	isAudio := utils.IsAudio(filename)
	contentType := utils.GetContentType(filename)

	fileType := "other"
	if isImage {
		fileType = "image"
	} else if isVideo {
		fileType = "video"
	} else if isAudio {
		fileType = "audio"
	}

	// Build URLs
	urls := map[string]string{
		"download": fmt.Sprintf("%s/api/files/%s", h.Config.BaseURL, filename),
		"view":     fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, filename),
		"metadata": fmt.Sprintf("%s/api/files/metadata/%s", h.Config.BaseURL, filename),
	}

	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	// If image, check for resized versions
	if isImage {
		resolutions := []string{"thumbnail", "small", "medium", "large"}

		for _, res := range resolutions {
			resizedFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, res, ext)
			resizedPath := filepath.Join(h.Config.UploadDir, resizedFilename)
			if _, err := os.Stat(resizedPath); err == nil {
				urls[fmt.Sprintf("view_%s", res)] = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, resizedFilename)
			}
		}
	}

	// If video, check for processed versions
	if isVideo {
		// Check for thumbnail
		thumbnailFilename := fmt.Sprintf("%s_thumbnail.jpg", nameWithoutExt)
		thumbnailPath := filepath.Join(h.Config.UploadDir, thumbnailFilename)
		if _, err := os.Stat(thumbnailPath); err == nil {
			urls["thumbnail"] = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, thumbnailFilename)
		}

		// Check for video resolutions
		resolutions := []string{"360p", "480p", "720p", "1080p"}
		for _, res := range resolutions {
			resFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, res, ext)
			resPath := filepath.Join(h.Config.UploadDir, resFilename)
			if _, err := os.Stat(resPath); err == nil {
				urls[fmt.Sprintf("view_%s", res)] = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, resFilename)
			}
		}
	}

	// If audio, check for different bitrates
	if isAudio {
		bitrates := []string{"low", "medium", "high"}
		for _, quality := range bitrates {
			audioFilename := fmt.Sprintf("%s_%s.mp3", nameWithoutExt, quality)
			audioPath := filepath.Join(h.Config.UploadDir, audioFilename)
			if _, err := os.Stat(audioPath); err == nil {
				urls[fmt.Sprintf("audio_%s", quality)] = fmt.Sprintf("%s/api/files/view/%s", h.Config.BaseURL, audioFilename)
			}
		}
	}

	metadata := models.FileMetadata{
		Success:     true,
		FileName:    filename,
		FileSize:    fileInfo.Size(),
		ContentType: contentType,
		FileType:    fileType,
		IsImage:     isImage,
		IsVideo:     isVideo,
		IsAudio:     isAudio,
		UploadedAt:  fileInfo.ModTime().Format("2006-01-02T15:04:05Z07:00"),
		URLs:        urls,
	}

	return c.Status(fiber.StatusOK).JSON(metadata)
}
