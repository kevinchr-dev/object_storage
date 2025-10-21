# Quick Start Guide

Panduan cepat untuk memulai menggunakan Object Storage Server.

## 🚀 Quick Installation (5 Menit)

### Step 1: Install FFmpeg

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt update && sudo apt install ffmpeg -y
```

**Verify:**
```bash
ffmpeg -version
```

### Step 2: Install Dependencies

```bash
cd /Users/kevinchr/object_storage
go mod download
```

### Step 3: Start Server

```bash
# Option 1: Run langsung
go run main.go

# Option 2: Build dulu
go build -o object-storage-server .
./object-storage-server
```

✅ Server running di: **http://localhost:8080**

---

## 📤 Quick Upload Test

### Via Browser (Paling Mudah)

1. Buka: `test-upload.html` di browser
2. Pilih file (gambar, video, atau audio)
3. Click "Upload File"
4. Lihat hasilnya!

### Via cURL

**Upload Image:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/image.jpg"
```

**Upload Video:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/video.mp4"
```

**Upload Audio:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/audio.mp3"
```

---

## 🎯 What Happens After Upload?

### Images (Instant Processing)
Langsung dapat 5 files:
- ✅ Original
- ✅ Large (1920px)
- ✅ Medium (1024px)
- ✅ Small (480px)
- ✅ Thumbnail (150px)

### Videos (Background Processing ~30-60s)
Akan dapat 6 files:
- ⏳ Original (instant)
- ⏳ 1080p HD (processing...)
- ⏳ 720p (processing...)
- ⏳ 480p (processing...)
- ⏳ 360p (processing...)
- ⏳ Thumbnail JPG (processing...)

### Audio (Background Processing ~10-30s)
Akan dapat 4 files:
- ⏳ Original (instant)
- ⏳ High 320k (processing...)
- ⏳ Medium 128k (processing...)
- ⏳ Low 64k (processing...)

### Other Files (PDF, DOCX, etc)
- ✅ Original only

---

## 📊 Check Processing Status

### Get Metadata (Recommended)

```bash
# Ganti dengan filename yang didapat dari upload response
curl http://localhost:8080/api/files/metadata/019a0566-fbb2-77a5-b1f8-43196337be36.mp4
```

**Response akan menunjukkan:**
- ✅ URLs yang sudah available
- ⏳ URLs yang masih processing (akan 404 jika belum selesai)

---

## 🎨 Use in Frontend

### Responsive Image (HTML)

```html
<img 
  src="http://localhost:8080/api/files/view/UUID_medium.jpg"
  srcset="
    http://localhost:8080/api/files/view/UUID_small.jpg 480w,
    http://localhost:8080/api/files/view/UUID_medium.jpg 1024w,
    http://localhost:8080/api/files/view/UUID_large.jpg 1920w
  "
  sizes="(max-width: 480px) 480px, (max-width: 1024px) 1024px, 1920px"
  alt="Responsive Image"
>
```

### Video Player with Quality Selection

```html
<video controls poster="http://localhost:8080/api/files/view/UUID_thumb.jpg">
  <source src="http://localhost:8080/api/files/view/UUID_1080p.mp4" label="1080p">
  <source src="http://localhost:8080/api/files/view/UUID_720p.mp4" label="720p">
  <source src="http://localhost:8080/api/files/view/UUID_480p.mp4" label="480p">
  <source src="http://localhost:8080/api/files/view/UUID_360p.mp4" label="360p">
</video>
```

### Audio Player with Bitrate Selection

```html
<audio controls>
  <source src="http://localhost:8080/api/files/view/UUID_high.mp3" type="audio/mpeg">
  <source src="http://localhost:8080/api/files/view/UUID_medium.mp3" type="audio/mpeg">
  <source src="http://localhost:8080/api/files/view/UUID_low.mp3" type="audio/mpeg">
</audio>
```

### JavaScript Fetch API

```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const response = await fetch('http://localhost:8080/api/upload', {
  method: 'POST',
  body: formData
});

const data = await response.json();
console.log('Upload success:', data);

// For video/audio, wait for processing
if (data.is_video || data.is_audio) {
  console.log('Processing in background...');
  
  // Poll metadata setiap 5 detik
  const interval = setInterval(async () => {
    const meta = await fetch(data.metadata_url).then(r => r.json());
    
    // Check if all URLs available
    if (meta.urls.view_720p) { // or any processed URL
      console.log('Processing complete!', meta.urls);
      clearInterval(interval);
    }
  }, 5000);
}
```

---

## ⚡ Performance Tips

### 1. Image Optimization
- Gunakan `thumbnail` untuk list/preview
- Gunakan `small` untuk mobile devices
- Gunakan `medium` untuk tablet/desktop
- Gunakan `large` hanya untuk full-screen viewing
- `original` hanya untuk download

### 2. Video Optimization
- Default ke `480p` untuk mobile
- Auto-detect connection speed untuk quality selection
- Preload `thumbnail` untuk instant preview
- Use `360p` untuk slow connections

### 3. Audio Optimization
- Start dengan `low` bitrate untuk instant playback
- Progressive upgrade to `medium` setelah buffering
- Use `high` hanya untuk premium users atau WiFi

---

## 🔧 Troubleshooting

### FFmpeg Not Found
```bash
# Pastikan FFmpeg terinstall
ffmpeg -version

# macOS install
brew install ffmpeg

# Ubuntu/Debian install
sudo apt install ffmpeg
```

### Video Processing Stuck
- Check server logs untuk error messages
- Pastikan FFmpeg terinstall dengan benar
- Video size terlalu besar? (default max: 50MB)

### Audio Quality Issues
- Original audio bitrate mungkin lebih rendah dari target
- FFmpeg akan maintain original bitrate jika lebih rendah

### Port Already in Use
```bash
# Change port di .env
PORT=3000

# Atau set environment variable
PORT=3000 go run main.go
```

---

## 📚 Next Steps

- 📖 Read full [README.md](README.md)
- 🎨 Check [IMAGE_PROCESSING.md](IMAGE_PROCESSING.md)
- 🎥 Check [VIDEO_AUDIO_PROCESSING.md](VIDEO_AUDIO_PROCESSING.md)
- 🔑 Learn about [UUID_V7.md](UUID_V7.md)
- 💻 See [API_EXAMPLES.md](API_EXAMPLES.md) for more code samples

---

## 🎯 Production Checklist

- [ ] Install FFmpeg on production server
- [ ] Set proper `MAX_FILE_SIZE` limit
- [ ] Configure CORS for your domain
- [ ] Set up reverse proxy (nginx/Apache)
- [ ] Enable HTTPS
- [ ] Set up CDN for static file serving
- [ ] Monitor disk space (processed files = more storage)
- [ ] Set up log rotation
- [ ] Configure file retention policy
- [ ] Set up backup for uploads directory

---

**Happy coding! 🚀**

For issues or questions, check the full documentation in [README.md](README.md)
