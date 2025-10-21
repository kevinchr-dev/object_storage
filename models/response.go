package models

type ViewURLs struct {
	Original  string `json:"original"`
	Large     string `json:"large,omitempty"`
	Medium    string `json:"medium,omitempty"`
	Small     string `json:"small,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	// Video resolutions
	Video1080p string `json:"1080p,omitempty"`
	Video720p  string `json:"720p,omitempty"`
	Video480p  string `json:"480p,omitempty"`
	Video360p  string `json:"360p,omitempty"`
	// Audio bitrates
	AudioHigh   string `json:"high,omitempty"`
	AudioMedium string `json:"medium,omitempty"`
	AudioLow    string `json:"low,omitempty"`
}

type UploadResponse struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	FileName    string    `json:"file_name,omitempty"`
	FileURL     string    `json:"file_url,omitempty"`
	ViewURLs    *ViewURLs `json:"view_urls,omitempty"`
	MetadataURL string    `json:"metadata_url,omitempty"`
	FileSize    int64     `json:"file_size,omitempty"`
	FileType    string    `json:"file_type,omitempty"` // "image", "video", "audio", "other"
	IsImage     bool      `json:"is_image,omitempty"`
	IsVideo     bool      `json:"is_video,omitempty"`
	IsAudio     bool      `json:"is_audio,omitempty"`
}

type FileMetadata struct {
	Success     bool              `json:"success"`
	FileName    string            `json:"file_name"`
	FileSize    int64             `json:"file_size"`
	ContentType string            `json:"content_type"`
	FileType    string            `json:"file_type"` // "image", "video", "audio", "other"
	IsImage     bool              `json:"is_image"`
	IsVideo     bool              `json:"is_video"`
	IsAudio     bool              `json:"is_audio"`
	UploadedAt  string            `json:"uploaded_at"`
	URLs        map[string]string `json:"urls"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
