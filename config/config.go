package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort   string
	UploadDir    string
	MaxFileSize  int64
	AllowedHosts string
	BaseURL      string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	maxFileSizeStr := os.Getenv("MAX_FILE_SIZE")
	maxFileSize := int64(4 * 1024 * 1024 * 1024) // Default 4GB
	if maxFileSizeStr != "" {
		if size, err := strconv.ParseInt(maxFileSizeStr, 10, 64); err == nil {
			maxFileSize = size
		}
	}

	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	if allowedHosts == "" {
		allowedHosts = "*"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	return &Config{
		ServerPort:   port,
		UploadDir:    uploadDir,
		MaxFileSize:  maxFileSize,
		AllowedHosts: allowedHosts,
		BaseURL:      baseURL,
	}
}
