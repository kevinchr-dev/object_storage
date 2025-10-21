package utils

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// GenerateUniqueFileName generates a unique filename using UUID v7
func GenerateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	uuidV7 := uuid.Must(uuid.NewV7()) // UUID v7 - time-ordered UUID
	return fmt.Sprintf("%s%s", uuidV7.String(), ext)
}

// IsImage checks if file is an image based on extension
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// IsVideo checks if file is a video based on extension
func IsVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	videoExts := []string{".mp4", ".avi", ".mov", ".mkv", ".webm", ".flv", ".wmv", ".m4v"}
	for _, vidExt := range videoExts {
		if ext == vidExt {
			return true
		}
	}
	return false
}

// IsAudio checks if file is an audio based on extension
func IsAudio(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	audioExts := []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma"}
	for _, audExt := range audioExts {
		if ext == audExt {
			return true
		}
	}
	return false
}

// CheckFFmpegInstalled checks if FFmpeg is installed
func CheckFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// ResizeImage creates multiple resized versions of an image
func ResizeImage(inputPath, outputDir, baseFilename string) (map[string]string, error) {
	// Open original image
	src, err := imaging.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	// Define resolutions
	resolutions := map[string]int{
		"thumbnail": 150,
		"small":     480,
		"medium":    1024,
		"large":     1920,
	}

	resizedFiles := make(map[string]string)
	ext := filepath.Ext(baseFilename)
	nameWithoutExt := strings.TrimSuffix(baseFilename, ext)

	for name, width := range resolutions {
		// Skip if original is smaller than target resolution
		bounds := src.Bounds()
		if bounds.Dx() <= width {
			continue
		}

		// Resize image maintaining aspect ratio
		resized := imaging.Resize(src, width, 0, imaging.Lanczos)

		// Generate filename
		resizedFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, name, ext)
		resizedPath := filepath.Join(outputDir, resizedFilename)

		// Save based on format
		if err := saveImage(resized, resizedPath, ext); err != nil {
			continue // Skip if save fails
		}

		resizedFiles[name] = resizedFilename
	}

	return resizedFiles, nil
}

// saveImage saves an image based on its extension
func saveImage(img image.Image, path, ext string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 85})
	case ".png":
		return png.Encode(out, img)
	case ".gif":
		return gif.Encode(out, img, nil)
	default:
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 85})
	}
}

// GetContentType returns content type based on file extension
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".txt":  "text/plain",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".mkv":  "video/x-matroska",
		".webm": "video/webm",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".flac": "audio/flac",
		".aac":  "audio/aac",
		".ogg":  "audio/ogg",
		".m4a":  "audio/mp4",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}

// ProcessVideo creates thumbnail and multiple resolutions for video
func ProcessVideo(inputPath, outputDir, baseFilename string) (map[string]string, error) {
	if !CheckFFmpegInstalled() {
		return nil, fmt.Errorf("ffmpeg not installed")
	}

	processedFiles := make(map[string]string)
	ext := filepath.Ext(baseFilename)
	nameWithoutExt := strings.TrimSuffix(baseFilename, ext)

	// Generate thumbnail (frame at 1 second)
	thumbnailFilename := fmt.Sprintf("%s_thumbnail.jpg", nameWithoutExt)
	thumbnailPath := filepath.Join(outputDir, thumbnailFilename)

	err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": "00:00:01"}).
		Output(thumbnailPath, ffmpeg.KwArgs{
			"vframes": 1,
			"vf":      "scale=320:-1",
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	if err == nil {
		processedFiles["thumbnail"] = thumbnailFilename
	}

	// Generate different resolutions
	resolutions := map[string]string{
		"360p":  "640:360",
		"480p":  "854:480",
		"720p":  "1280:720",
		"1080p": "1920:1080",
	}

	for quality, scale := range resolutions {
		resFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, quality, ext)
		resPath := filepath.Join(outputDir, resFilename)

		err := ffmpeg.Input(inputPath).
			Output(resPath, ffmpeg.KwArgs{
				"vf":     fmt.Sprintf("scale=%s", scale),
				"c:v":    "libx264",
				"crf":    "23",
				"c:a":    "aac",
				"b:a":    "128k",
				"preset": "fast",
			}).
			OverWriteOutput().
			ErrorToStdOut().
			Run()

		if err == nil {
			processedFiles[quality] = resFilename
		}
	}

	return processedFiles, nil
}

// ProcessAudio creates multiple bitrates for audio
func ProcessAudio(inputPath, outputDir, baseFilename string) (map[string]string, error) {
	if !CheckFFmpegInstalled() {
		return nil, fmt.Errorf("ffmpeg not installed")
	}

	processedFiles := make(map[string]string)
	ext := filepath.Ext(baseFilename)
	nameWithoutExt := strings.TrimSuffix(baseFilename, ext)

	// Define bitrates
	bitrates := map[string]string{
		"low":    "64k",
		"medium": "128k",
		"high":   "320k",
	}

	for quality, bitrate := range bitrates {
		audioFilename := fmt.Sprintf("%s_%s.mp3", nameWithoutExt, quality)
		audioPath := filepath.Join(outputDir, audioFilename)

		err := ffmpeg.Input(inputPath).
			Output(audioPath, ffmpeg.KwArgs{
				"b:a": bitrate,
				"c:a": "libmp3lame",
				"ar":  "44100",
			}).
			OverWriteOutput().
			ErrorToStdOut().
			Run()

		if err == nil {
			processedFiles[quality] = audioFilename
		}
	}

	return processedFiles, nil
}
