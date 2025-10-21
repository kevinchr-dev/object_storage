# Docker Deployment Guide

Panduan lengkap untuk deploy Object Storage Server menggunakan Docker.

## 📋 Prerequisites

- Docker Engine 20.10+
- Docker Compose v2.0+
- Minimal 2GB RAM
- Minimal 10GB disk space

## 🚀 Quick Start (Development)

### 1. Clone & Setup

```bash
cd /path/to/object_storage
cp .env.docker .env
```

### 2. Build & Run

```bash
# Build dan start container
docker-compose up -d

# Check logs
docker-compose logs -f

# Check health
curl http://localhost:8080/api/health
```

### 3. Test Upload

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@test.jpg"
```

Server running di: **http://localhost:8080**

---

## 🏭 Production Deployment

### 1. Setup Environment

```bash
# Copy dan edit environment variables
cp .env.docker .env
nano .env
```

Edit `.env`:
```env
PORT=8080
BASE_URL=https://yourdomain.com
MAX_FILE_SIZE=104857600
CORS_ALLOW_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### 2. Deploy dengan Nginx Reverse Proxy

```bash
# Deploy dengan production config
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps
```

### 3. SSL Configuration (Optional)

Jika menggunakan SSL, tambahkan certificate:

```bash
# Create SSL directory
mkdir -p nginx/ssl

# Copy SSL certificates
cp /path/to/cert.pem nginx/ssl/
cp /path/to/key.pem nginx/ssl/

# Edit nginx.conf dan uncomment HTTPS server block
nano nginx/nginx.conf
```

---

## 🐳 Docker Commands

### Build & Run

```bash
# Build image
docker-compose build

# Start services
docker-compose up -d

# Start dengan rebuild
docker-compose up -d --build

# Stop services
docker-compose down

# Stop dan remove volumes
docker-compose down -v
```

### Logs & Monitoring

```bash
# View all logs
docker-compose logs

# Follow logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f object-storage

# Check resource usage
docker stats
```

### Container Management

```bash
# List running containers
docker-compose ps

# Restart service
docker-compose restart object-storage

# Execute command in container
docker-compose exec object-storage sh

# View container details
docker inspect object-storage-server
```

---

## 📊 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `BASE_URL` | http://localhost:8080 | Base URL untuk file URLs |
| `UPLOAD_DIR` | /app/uploads | Upload directory path |
| `MAX_FILE_SIZE` | 52428800 | Max file size (50MB) |
| `CORS_ALLOW_ORIGINS` | * | CORS allowed origins |

### Volume Mounts

```yaml
volumes:
  - ./uploads:/app/uploads        # Persistent file storage
  - ./.env:/app/.env:ro           # Environment config
```

**Important:** Directory `./uploads` akan dibuat otomatis dan berisi semua uploaded files.

### Resource Limits

Default limits di docker-compose.yml:

```yaml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 2G
    reservations:
      cpus: '0.5'
      memory: 512M
```

Adjust sesuai kebutuhan server Anda.

---

## 🔧 Customization

### Increase Max File Size

Edit `docker-compose.yml`:

```yaml
environment:
  - MAX_FILE_SIZE=104857600  # 100MB
```

Edit `nginx/nginx.conf`:

```nginx
client_max_body_size 100M;
```

### Add Multiple Instances (Load Balancing)

Edit `docker-compose.yml`:

```yaml
services:
  object-storage:
    # ... existing config
    deploy:
      replicas: 3
```

---

## 🚨 Troubleshooting

### Container won't start

```bash
# Check logs for errors
docker-compose logs object-storage

# Check if port is already in use
lsof -i :8080

# Remove old containers and try again
docker-compose down
docker-compose up -d
```

### FFmpeg not working

FFmpeg sudah included di Docker image. Jika ada masalah:

```bash
# Verify FFmpeg in container
docker-compose exec object-storage ffmpeg -version
```

### Upload directory permissions

```bash
# Fix permissions
sudo chown -R 1000:1000 ./uploads
```

### Out of disk space

```bash
# Check disk usage
df -h

# Clean old Docker images
docker system prune -a

# Remove old uploaded files
rm -rf ./uploads/*
```

---

## 📈 Performance Tuning

### For High Traffic

Edit `docker-compose.yml`:

```yaml
deploy:
  resources:
    limits:
      cpus: '4'
      memory: 4G
    reservations:
      cpus: '2'
      memory: 2G
```

Edit `nginx/nginx.conf`:

```nginx
worker_processes auto;
worker_connections 8192;
```

### Enable Caching

Nginx sudah dikonfigurasi dengan caching. Untuk customize:

```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=10g inactive=60m;
```

---

## 🔐 Security Best Practices

### 1. Jangan expose port langsung

Gunakan nginx reverse proxy:

```yaml
# object-storage service
expose:
  - "8080"  # Not 'ports' - only nginx can access
```

### 2. Set CORS dengan benar

```env
CORS_ALLOW_ORIGINS=https://yourdomain.com
```

### 3. Enable HTTPS

Uncomment HTTPS block di `nginx/nginx.conf` dan tambahkan SSL certificates.

### 4. Rate Limiting

Sudah dikonfigurasi di nginx:
- Upload: 10 req/s per IP
- Download: 50 req/s per IP

### 5. Regular Updates

```bash
# Update base images
docker-compose pull
docker-compose up -d
```

---

## 📦 Backup & Restore

### Backup Uploaded Files

```bash
# Backup uploads directory
tar -czf uploads-backup-$(date +%Y%m%d).tar.gz uploads/

# Or use rsync for incremental backup
rsync -av ./uploads/ /backup/location/
```

### Restore

```bash
# Extract backup
tar -xzf uploads-backup-20241021.tar.gz

# Restart service
docker-compose restart object-storage
```

---

## 🌐 Production Checklist

- [ ] Set proper `BASE_URL` in `.env`
- [ ] Configure CORS origins (tidak `*`)
- [ ] Setup SSL certificates
- [ ] Enable HTTPS in nginx config
- [ ] Set resource limits sesuai server
- [ ] Setup backup automation
- [ ] Configure monitoring (Prometheus/Grafana)
- [ ] Setup log rotation
- [ ] Test health check endpoint
- [ ] Test file upload & download
- [ ] Load testing
- [ ] Setup firewall rules

---

## 📚 Additional Resources

- **[README.md](README.md)** - General documentation
- **[QUICK_START.md](QUICK_START.md)** - Quick start guide
- **Dockerfile** - Container build configuration
- **docker-compose.yml** - Development setup
- **docker-compose.prod.yml** - Production setup with nginx

---

## 🆘 Support

Jika ada masalah:

1. Check logs: `docker-compose logs -f`
2. Check health: `curl http://localhost:8080/api/health`
3. Check container status: `docker-compose ps`
4. Check resource usage: `docker stats`

---

**Last Updated:** October 21, 2024
