# 4GB File Size Limit Configuration

## Overview
Server ini mendukung upload dan download file hingga **4GB per request**. Untuk file yang lebih besar dari 4GB, gunakan multipart upload/download dari sisi client.

## Konfigurasi yang Diterapkan

### 1. Go Application (config/config.go)
```go
// Default 4GB
maxFileSize := int64(4 * 1024 * 1024 * 1024)
```

### 2. Fiber Framework (main.go)
```go
app := fiber.New(fiber.Config{
    BodyLimit:        int(cfg.MaxFileSize), // 4GB
    StreamRequestBody: true,                 // Stream untuk efisiensi memory
    ReadBufferSize:   16384,                // 16KB buffer
    WriteBufferSize:  16384,                // 16KB buffer
    ReadTimeout:      0,                    // No timeout untuk file besar
    WriteTimeout:     0,                    // No timeout untuk file besar
})
```

### 3. Docker Environment (.env.docker)
```bash
MAX_FILE_SIZE=4294967296  # 4GB dalam bytes
```

### 4. Nginx Proxy (nginx/nginx.conf)
```nginx
# Maximum upload size
client_max_body_size 4G;

# Timeouts untuk upload/download file besar (1 jam)
proxy_connect_timeout 600s;
proxy_send_timeout 3600s;
proxy_read_timeout 3600s;
```

## Testing Upload 4GB File

### 1. Test dengan cURL
```bash
# Upload file 4GB
curl -X POST http://localhost:8080/api/upload \
  -F "file=@large_file.mp4" \
  -w "\nTime: %{time_total}s\n"
```

### 2. Test dengan Script Bash
```bash
#!/bin/bash
# create_and_upload_large_file.sh

# Buat dummy file 4GB (opsional - jika tidak punya file test)
# dd if=/dev/zero of=test_4gb.bin bs=1M count=4096

# Upload
echo "Uploading 4GB file..."
time curl -X POST http://localhost:8080/api/upload \
  -F "file=@test_4gb.bin" \
  -o upload_response.json

cat upload_response.json | jq '.'
```

### 3. Test Download
```bash
# Download file yang sudah diupload
FILE_ID="your-uuid-v7-here"
time curl -X GET "http://localhost:8080/api/files/${FILE_ID}/download" \
  -o downloaded_file.bin \
  -w "\nTime: %{time_total}s\nSize: %{size_download} bytes\n"
```

## Multipart Upload untuk File >4GB (Client-Side)

Untuk file lebih besar dari 4GB, implementasikan multipart upload di sisi client:

### JavaScript/Browser Example
```javascript
// multipart_upload.js
async function uploadLargeFile(file) {
    const CHUNK_SIZE = 4 * 1024 * 1024 * 1024; // 4GB per chunk
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    const uploadedParts = [];

    for (let i = 0; i < totalChunks; i++) {
        const start = i * CHUNK_SIZE;
        const end = Math.min(start + CHUNK_SIZE, file.size);
        const chunk = file.slice(start, end);

        const formData = new FormData();
        formData.append('file', chunk, `${file.name}.part${i}`);

        const response = await fetch('/api/upload', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();
        uploadedParts.push(result.data.file_id);
        
        console.log(`Uploaded part ${i + 1}/${totalChunks}`);
    }

    // Merge parts di server (perlu endpoint tambahan)
    const mergeResponse = await fetch('/api/merge', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            parts: uploadedParts,
            filename: file.name
        })
    });

    return await mergeResponse.json();
}

// Usage
const fileInput = document.getElementById('file-input');
fileInput.addEventListener('change', async (e) => {
    const file = e.target.files[0];
    if (file.size > 4 * 1024 * 1024 * 1024) {
        console.log('File >4GB, using multipart upload...');
        await uploadLargeFile(file);
    } else {
        // Regular upload
        const formData = new FormData();
        formData.append('file', file);
        await fetch('/api/upload', { method: 'POST', body: formData });
    }
});
```

### Python Example
```python
# multipart_upload.py
import requests
import os
import math

def upload_large_file(filepath):
    """Upload file dengan chunking 4GB per request"""
    CHUNK_SIZE = 4 * 1024 * 1024 * 1024  # 4GB
    file_size = os.path.getsize(filepath)
    total_chunks = math.ceil(file_size / CHUNK_SIZE)
    uploaded_parts = []
    
    with open(filepath, 'rb') as f:
        for i in range(total_chunks):
            chunk = f.read(CHUNK_SIZE)
            
            files = {
                'file': (f'{os.path.basename(filepath)}.part{i}', chunk)
            }
            
            response = requests.post(
                'http://localhost:8080/api/upload',
                files=files
            )
            
            result = response.json()
            uploaded_parts.append(result['data']['file_id'])
            
            print(f"Uploaded part {i + 1}/{total_chunks}")
    
    # Merge parts (perlu endpoint tambahan)
    merge_response = requests.post(
        'http://localhost:8080/api/merge',
        json={
            'parts': uploaded_parts,
            'filename': os.path.basename(filepath)
        }
    )
    
    return merge_response.json()

if __name__ == '__main__':
    # Upload file >4GB
    result = upload_large_file('/path/to/large_file.mp4')
    print(f"Upload complete: {result}")
```

### Go Client Example
```go
// multipart_upload.go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)

const ChunkSize = 4 * 1024 * 1024 * 1024 // 4GB

type UploadResponse struct {
    Success bool   `json:"success"`
    Data    struct {
        FileID string `json:"file_id"`
    } `json:"data"`
}

func uploadLargeFile(filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    fileInfo, _ := file.Stat()
    totalChunks := (fileInfo.Size() + ChunkSize - 1) / ChunkSize
    var uploadedParts []string

    buffer := make([]byte, ChunkSize)
    
    for i := int64(0); i < totalChunks; i++ {
        n, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            return err
        }

        // Upload chunk
        var b bytes.Buffer
        w := multipart.NewWriter(&b)
        
        fw, _ := w.CreateFormFile("file", fmt.Sprintf("%s.part%d", 
            filepath.Base(filepath), i))
        fw.Write(buffer[:n])
        w.Close()

        resp, err := http.Post("http://localhost:8080/api/upload", 
            w.FormDataContentType(), &b)
        if err != nil {
            return err
        }

        var result UploadResponse
        json.NewDecoder(resp.Body).Decode(&result)
        resp.Body.Close()

        uploadedParts = append(uploadedParts, result.Data.FileID)
        
        fmt.Printf("Uploaded part %d/%d\n", i+1, totalChunks)
    }

    // Merge parts (perlu endpoint tambahan)
    mergeData, _ := json.Marshal(map[string]interface{}{
        "parts": uploadedParts,
        "filename": filepath.Base(filepath),
    })

    resp, err := http.Post("http://localhost:8080/api/merge",
        "application/json", bytes.NewBuffer(mergeData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    fmt.Println("Upload complete!")
    return nil
}

func main() {
    uploadLargeFile("/path/to/large_file.mp4")
}
```

## Performance Considerations

### 1. Upload Speed Estimation
- **Network 1 Gbps**: ~32 seconds untuk 4GB
- **Network 100 Mbps**: ~5.3 menit untuk 4GB
- **Network 10 Mbps**: ~53 menit untuk 4GB

### 2. Memory Usage
- Streaming mode: ~16KB buffer (minimal memory)
- Worker pool: 4 concurrent processing jobs
- Nginx cache: Configurable (default 10GB)

### 3. Disk Space
- Raw file: 4GB
- Image resolutions: +20-100MB (4 sizes)
- Video resolutions: +500MB-2GB (4 sizes + thumbnail)
- Audio bitrates: +20-80MB (3 bitrates)
- **Total per 4GB video**: ~6-7GB storage

### 4. Timeout Settings
```nginx
# Upload endpoint
proxy_connect_timeout 600s;   # 10 menit untuk koneksi
proxy_send_timeout 3600s;      # 1 jam untuk upload
proxy_read_timeout 3600s;      # 1 jam untuk response

# Download endpoint
proxy_send_timeout 3600s;      # 1 jam untuk download
proxy_read_timeout 3600s;      # 1 jam untuk read
```

## Monitoring & Troubleshooting

### Check Upload Progress
```bash
# Monitor nginx logs
docker-compose logs -f nginx

# Monitor Go app logs
docker-compose logs -f object-storage
```

### Check Disk Usage
```bash
# Check uploads directory
du -sh uploads/

# Check detailed breakdown
du -h --max-depth=1 uploads/
```

### Common Issues

1. **413 Request Entity Too Large**
   - Check nginx `client_max_body_size`
   - Check Fiber `BodyLimit`
   - Verify environment variable `MAX_FILE_SIZE`

2. **504 Gateway Timeout**
   - Increase nginx `proxy_read_timeout`
   - Increase nginx `proxy_send_timeout`
   - Check network speed

3. **Out of Memory**
   - Verify `StreamRequestBody: true` di Fiber
   - Check buffer sizes (should be 16KB)
   - Monitor `docker stats`

4. **Slow Upload**
   - Check network bandwidth
   - Disable compression untuk file besar
   - Consider multipart untuk file >4GB

## Recommendations

1. **Untuk File <4GB**: Upload langsung dengan single request
2. **Untuk File >4GB**: Gunakan multipart upload client-side
3. **Video Processing**: Akan memakan waktu lama, gunakan worker pool (async)
4. **Monitor**: Setup monitoring untuk disk space dan memory usage
5. **Backup**: Implement backup strategy untuk files penting

## Next Steps

Jika Anda perlu menambahkan endpoint untuk merge multipart files:

1. Tambahkan endpoint `/api/merge` di `handlers/file_handler.go`
2. Implementasi fungsi untuk menggabungkan chunk files
3. Update worker pool untuk handle merge jobs
4. Tambahkan validasi untuk memastikan semua parts tersedia
5. Cleanup temporary part files setelah merge sukses

Apakah Anda perlu saya buatkan endpoint merge untuk multipart upload?
