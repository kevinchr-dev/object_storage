# Object Storage Server

API service untuk object storage server yang dibangun dengan Golang dan Go Fiber framework. Server ini memungkinkan client untuk upload file, mendapatkan URL, dan melakukan view/download file.

## Fitur

### 📤 Upload & Storage
- ✅ Upload file (semua jenis file)
- ✅ **UUID v7** untuk nama file (secure & sortable)
- ✅ **Streaming upload** - support large files tanpa memory overflow
- ✅ View file secara inline di browser
- ✅ Download file dengan streaming
- ✅ Secure filename generation

### 🎨 Media Processing
- ✅ **Multiple resolutions** untuk gambar (thumbnail, small, medium, large)
- ✅ **Video processing** dengan multiple resolutions (360p, 480p, 720p, 1080p) + thumbnail
- ✅ **Audio processing** dengan multiple bitrates (64k, 128k, 320k)
- ✅ **Auto image compression** dengan quality optimization
- ✅ **Smart processing** - instant untuk file kecil, background untuk file besar

### ⚡ Performance & Scalability
- ✅ **256K concurrent connections** support
- ✅ **Worker Pool** dengan 4 concurrent processors
- ✅ **Non-blocking uploads** - instant response
- ✅ **Background processing** untuk video & audio
- ✅ **Rate limiting** - 100 req/min per IP
- ✅ **Response compression** - gzip untuk bandwidth optimization
- ✅ **Parallel processing** - handle multiple users simultaneously

### 🔧 API & Integration
- ✅ **Metadata endpoint** dengan informasi lengkap file
- ✅ CORS support
- ✅ RESTful API design
- ✅ File size validation
- ✅ Comprehensive error handling

### 🐳 Deployment
- ✅ **Docker support** - containerized deployment
- ✅ **Docker Compose** - easy multi-container setup
- ✅ **Production ready** - nginx reverse proxy included
- ✅ **One-command deploy** - `docker-compose up -d`

## Teknologi

- **Golang** - Programming language
- **Fiber v2** - Web framework
- **UUID v7** - Unique filename generation (time-ordered)
- **Imaging Library** - Image processing & resizing
- **MD5** - Hash generation untuk verification

## 🚀 Quick Start

### Option 1: Docker (Recommended)

```bash
# 1. Clone repository
git clone <repository-url>
cd object_storage

# 2. Start with Docker Compose
docker-compose up -d

# 3. Test
curl http://localhost:8080/api/health
```

**That's it!** Server running di `http://localhost:8080` ✅

📚 **[Full Docker Guide →](DOCKER.md)**

---

### Option 2: Manual Installation

#### Prerequisites

- **Go 1.16 atau lebih tinggi**
- **FFmpeg** (diperlukan untuk video & audio processing)
- Git

#### Install FFmpeg

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
```

**Verify instalasi:**
```bash
ffmpeg -version
```

#### Setup

1. Clone repository atau buka folder project

2. Install dependencies:
```bash
go mod tidy
```

3. Copy file environment:
```bash
cp .env.example .env
```

4. Sesuaikan konfigurasi di file `.env` jika diperlukan:
```env
PORT=3000
BASE_URL=http://localhost:3000
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=52428800
ALLOWED_HOSTS=*
```

5. Jalankan server:
```bash
go run main.go
```

Server akan berjalan di `http://localhost:3000`

## API Endpoints

### 1. Upload File

**POST** `/api/upload`

Upload file ke server. File akan otomatis mendapat nama UUID v7. Jika file adalah gambar, akan otomatis di-generate multiple resolusi.

**Request:**
- Method: POST
- Content-Type: multipart/form-data
- Body: 
  - `file`: File yang akan diupload

**Example (cURL):**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/your/file.jpg"
```

**Response (Non-Image):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "file_name": "550e8400-e29b-41d4-a716-446655440000.pdf",
  "file_url": "http://localhost:8080/api/files/550e8400-e29b-41d4-a716-446655440000.pdf",
  "view_urls": {
    "original": "http://localhost:8080/api/files/view/550e8400-e29b-41d4-a716-446655440000.pdf"
  },
  "metadata_url": "http://localhost:8080/api/files/metadata/550e8400-e29b-41d4-a716-446655440000.pdf",
  "file_size": 245632,
  "is_image": false
}
```

**Response (Image - dengan multiple resolusi):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "file_name": "019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
  "file_url": "http://localhost:8080/api/files/019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
  "view_urls": {
    "original": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
    "large": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_large.jpg",
    "medium": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_medium.jpg",
    "small": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_small.jpg",
    "thumbnail": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_thumbnail.jpg"
  },
  "metadata_url": "http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
  "file_size": 2456320,
  "is_image": true
}
```

**Response (Video - dengan background processing):**
```json
{
  "success": true,
  "message": "File uploaded successfully. Video processing started in background.",
  "file_name": "019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
  "file_url": "http://localhost:8080/api/files/019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
  "view_urls": {
    "original": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
    "video_1080p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_1080p.mp4",
    "video_720p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_720p.mp4",
    "video_480p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_480p.mp4",
    "video_360p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_360p.mp4",
    "thumbnail": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_thumb.jpg"
  },
  "metadata_url": "http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
  "file_size": 15728640,
  "is_video": true
}
```

**Response (Audio - dengan background processing):**
```json
{
  "success": true,
  "message": "File uploaded successfully. Audio processing started in background.",
  "file_name": "019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
  "file_url": "http://localhost:8080/api/files/019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
  "view_urls": {
    "original": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
    "audio_high": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38_high.mp3",
    "audio_medium": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38_medium.mp3",
    "audio_low": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38_low.mp3"
  },
  "metadata_url": "http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
  "file_size": 5242880,
  "is_audio": true
}
```

**File Processing:**
- **Images**: Processed instantly (4 resolutions: thumbnail 150px, small 480px, medium 1024px, large 1920px)
- **Videos**: Background processing (~30-60 seconds) - creates thumbnail + 4 resolutions (360p, 480p, 720p, 1080p)
- **Audio**: Background processing (~10-30 seconds) - creates 3 bitrates (64k low, 128k medium, 320k high)
- **Other files**: Original only

### 2. Download File

**GET** `/api/files/:filename`

Download file dengan nama file yang diberikan.

**Example:**
```bash
curl -O http://localhost:3000/api/files/1729512345678_a1b2c3d4.jpg
```

atau buka di browser:
```
http://localhost:3000/api/files/1729512345678_a1b2c3d4.jpg
```

### 3. View File

**GET** `/api/files/view/:filename`

View file secara inline di browser (tidak download).

**Example:**
```
http://localhost:3000/api/files/view/1729512345678_a1b2c3d4.jpg
```

### 4. Get File Metadata (Recommended)

**GET** `/api/files/metadata/:filename`

Mendapatkan metadata lengkap file termasuk semua URL yang tersedia (untuk gambar, video, atau audio).

**Example:**
```
http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be36.jpg
```

**Response (Image):**
```json
{
  "success": true,
  "file_name": "019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
  "file_size": 2456320,
  "content_type": "image/jpeg",
  "is_image": true,
  "uploaded_at": "2024-10-21T10:30:45+07:00",
  "urls": {
    "download": "http://localhost:8080/api/files/019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
    "view": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
    "metadata": "http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be36.jpg",
    "view_large": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_large.jpg",
    "view_medium": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_medium.jpg",
    "view_small": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_small.jpg",
    "view_thumbnail": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_thumbnail.jpg"
  }
}
```

**Response (Video - setelah processing selesai):**
```json
{
  "success": true,
  "file_name": "019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
  "file_size": 15728640,
  "content_type": "video/mp4",
  "is_video": true,
  "uploaded_at": "2024-10-21T10:30:45+07:00",
  "urls": {
    "download": "http://localhost:8080/api/files/019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
    "view": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
    "metadata": "http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be37.mp4",
    "view_1080p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_1080p.mp4",
    "view_720p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_720p.mp4",
    "view_480p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_480p.mp4",
    "view_360p": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_360p.mp4",
    "view_thumbnail": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be37_thumb.jpg"
  }
}
```

**Response (Audio - setelah processing selesai):**
```json
{
  "success": true,
  "file_name": "019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
  "file_size": 5242880,
  "content_type": "audio/mpeg",
  "is_audio": true,
  "uploaded_at": "2024-10-21T10:30:45+07:00",
  "urls": {
    "download": "http://localhost:8080/api/files/019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
    "view": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
    "metadata": "http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be38.mp3",
    "view_high": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38_high.mp3",
    "view_medium": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38_medium.mp3",
    "view_low": "http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be38_low.mp3"
  }
}
```

**Note**: Untuk video dan audio, processed files akan tersedia setelah background processing selesai (~30-60 detik untuk video, ~10-30 detik untuk audio).

### 5. Get File Info (Deprecated)

**GET** `/api/files/info/:filename`

⚠️ **Deprecated**: Gunakan `/api/files/metadata/:filename` untuk informasi lebih lengkap.

**Response:**
```json
{
  "success": true,
  "file_name": "550e8400-e29b-41d4-a716-446655440000.jpg",
  "file_size": 245632,
  "modified": "2024-10-21T10:30:45Z"
}
```

### 6. Health Check

**GET** `/api/health`

Check status server.

**Response:**
```json
{
  "status": "ok",
  "message": "Object Storage Server is running"
}
```

## Use Cases & Best Practices

### 1. Responsive Images (Bandwidth Optimization)

Gunakan resolusi yang sesuai untuk device yang berbeda:

```html
<!-- Responsive image dengan srcset -->
<img 
  src="http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_medium.jpg"
  srcset="
    http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_small.jpg 480w,
    http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_medium.jpg 1024w,
    http://localhost:8080/api/files/view/019a0566-fbb2-77a5-b1f8-43196337be36_large.jpg 1920w
  "
  sizes="(max-width: 480px) 480px, (max-width: 1024px) 1024px, 1920px"
  alt="Responsive Image"
>
```

### 2. Adaptive Video Streaming (Quality Selection)

Biarkan user memilih kualitas video berdasarkan bandwidth:

```html
<video controls>
  <source src="http://localhost:8080/api/files/view/UUID_1080p.mp4" label="1080p HD">
  <source src="http://localhost:8080/api/files/view/UUID_720p.mp4" label="720p">
  <source src="http://localhost:8080/api/files/view/UUID_480p.mp4" label="480p">
  <source src="http://localhost:8080/api/files/view/UUID_360p.mp4" label="360p">
</video>
```

**Auto Quality Selection:**
```javascript
const connection = navigator.connection;
let videoQuality = '720p'; // default

if (connection) {
  const effectiveType = connection.effectiveType;
  if (effectiveType === '4g') videoQuality = '1080p';
  else if (effectiveType === '3g') videoQuality = '480p';
  else if (effectiveType === '2g') videoQuality = '360p';
}

videoElement.src = `${baseUrl}_${videoQuality}.mp4`;
```

### 3. Progressive Audio Loading (Podcast/Music Apps)

Mulai dengan bitrate rendah, upgrade setelah buffering:

```javascript
const audioElement = new Audio();
audioElement.src = metadata.urls.view_low; // Start dengan 64k

audioElement.addEventListener('canplaythrough', () => {
  // Setelah low quality loaded, preload medium quality
  const mediumAudio = new Audio(metadata.urls.view_medium);
  mediumAudio.preload = 'auto';
  
  // Switch ke medium setelah 5 detik
  setTimeout(() => {
    audioElement.src = metadata.urls.view_medium;
    audioElement.play();
  }, 5000);
});
```

### 4. Thumbnail Preview Strategy

```javascript
// Load thumbnail dulu untuk fast preview
<img src="thumbnail_url" data-full="medium_url" class="lazy-load">

// Lazy load full resolution
observer.observe(img);
img.onload = () => {
  img.src = img.dataset.full;
};
```

### 5. CDN Integration

Metadata URL menyediakan semua informasi yang diperlukan untuk integrasi dengan CDN:

```javascript
const metadata = await fetch(metadataURL).then(r => r.json());
// Push semua URLs ke CDN
metadata.urls // contains all available URLs
```

## Contoh Penggunaan

### Upload dengan JavaScript (Fetch API)

```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

fetch('http://localhost:3000/api/upload', {
  method: 'POST',
  body: formData
})
  .then(response => response.json())
  .then(data => {
    console.log('File uploaded:', data.file_url);
    console.log('View URL:', data.view_url);
  })
  .catch(error => console.error('Error:', error));
```

### Upload dengan Python

```python
import requests

url = "http://localhost:3000/api/upload"
files = {'file': open('image.jpg', 'rb')}

response = requests.post(url, files=files)
print(response.json())
```

### Upload dengan Postman

1. Method: POST
2. URL: `http://localhost:3000/api/upload`
3. Body: form-data
4. Key: `file` (tipe: File)
5. Value: Pilih file yang ingin diupload

## Konfigurasi

### Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| PORT | 3000 | Port server |
| BASE_URL | http://localhost:3000 | Base URL untuk generate file URLs |
| UPLOAD_DIR | ./uploads | Directory untuk menyimpan file |
| MAX_FILE_SIZE | 52428800 | Maximum file size dalam bytes (default: 50MB) |
| ALLOWED_HOSTS | * | CORS allowed hosts |

## Struktur Project

```
object_storage/
├── config/
│   └── config.go          # Configuration management
├── handlers/
│   └── file_handler.go    # File upload/download handlers
├── models/
│   └── response.go        # Response models
├── routes/
│   └── routes.go          # API routes
├── utils/
│   └── file.go            # File utilities (UUID v7, image processing)
├── uploads/               # Upload directory (created automatically)
├── main.go                # Application entry point
├── go.mod                 # Go modules
├── .env.example           # Environment variables example
├── README.md              # Documentation (this file)
├── API_EXAMPLES.md        # API usage examples in multiple languages
├── IMAGE_PROCESSING.md    # Image processing & resolutions guide
└── test-upload.html       # Test UI for file upload
```

## 📚 Dokumentasi Lengkap

### Getting Started
- **[README.md](README.md)** - Setup & general documentation (you are here)
- **[QUICK_START.md](QUICK_START.md)** - ⚡ Quick start guide (5 menit setup!)
- **[DOCKER.md](DOCKER.md)** - 🐳 **Docker Deployment Guide (NEW!)**
- **[DOCKER_QUICK.md](DOCKER_QUICK.md)** - 🚀 **Quick Docker Commands**

### Performance & Architecture
- **[CONCURRENT_PERFORMANCE.md](CONCURRENT_PERFORMANCE.md)** - Concurrent & Performance Guide
- **[API_EXAMPLES.md](API_EXAMPLES.md)** - API usage dengan 10+ bahasa pemrograman
- **[IMAGE_PROCESSING.md](IMAGE_PROCESSING.md)** - Image processing & multiple resolutions guide
- **[VIDEO_AUDIO_PROCESSING.md](VIDEO_AUDIO_PROCESSING.md)** - Video & audio processing guide
- **[UUID_V7.md](UUID_V7.md)** - UUID v7 implementation & benefits
- **[test-upload.html](test-upload.html)** - Interactive test UI

## Quick Links

- 🐳 **[Docker Deploy Guide](DOCKER.md)** ← **Deploy dalam 2 menit!**
- ⚡ [Quick Start Guide](QUICK_START.md)
- 🚀 [Concurrent & Performance Guide](CONCURRENT_PERFORMANCE.md)
- 📖 [API Endpoints](#api-endpoints)
- 💡 [Use Cases & Best Practices](#use-cases--best-practices)
- 🔒 [Security Features](#security-features)

## Security Features

- **Directory Traversal Prevention**: Menggunakan `filepath.Base()` untuk prevent path traversal attacks
- **File Size Validation**: Validasi ukuran file sebelum upload
- **UUID v7 Filename**: Generate unique, time-ordered filename yang secure dan sortable
- **Content Type Detection**: Automatic content type detection berdasarkan file extension
- **Image Validation**: Validate image format sebelum processing

## Production Deployment

### Build Binary

```bash
go build -o object-storage-server
```

### Run Binary

```bash
./object-storage-server
```

### Dengan Docker (Optional)

Buat `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
RUN mkdir uploads
EXPOSE 3000
CMD ["./main"]
```

Build dan run:
```bash
docker build -t object-storage-server .
docker run -p 3000:3000 -v $(pwd)/uploads:/root/uploads object-storage-server
```

## Troubleshooting

### Port sudah digunakan
Ubah PORT di file `.env` ke port lain yang available.

### Permission denied saat create upload directory
Pastikan aplikasi memiliki write permission di directory tempat aplikasi dijalankan.

### File terlalu besar
Sesuaikan `MAX_FILE_SIZE` di file `.env` sesuai kebutuhan.

## License

MIT License

## Author

Kevin Christopher

## Support

Untuk pertanyaan dan dukungan, silakan buat issue di repository ini.
