# Quick Docker Deploy Commands

## 🚀 Quick Start

```bash
# 1. Build & Start
docker-compose up -d

# 2. Check status
docker-compose ps

# 3. View logs
docker-compose logs -f

# 4. Stop
docker-compose down
```

## 📝 Common Commands

```bash
# Start services
make up              # or: docker-compose up -d

# Stop services
make down            # or: docker-compose down

# Restart
make restart         # or: docker-compose restart

# View logs
make logs-f          # or: docker-compose logs -f

# Check health
make test            # or: curl http://localhost:8080/api/health

# Access shell
make shell           # or: docker-compose exec object-storage sh
```

## 🏭 Production Deploy

```bash
# With nginx reverse proxy
docker-compose -f docker-compose.prod.yml up -d

# Or using Make
make prod-up
```

## 🔧 Maintenance

```bash
# Backup uploads
make backup

# Clean containers & volumes
make clean

# Clean Docker system
make prune

# View resource usage
make stats
```

## 📊 Monitoring

```bash
# Real-time logs
docker-compose logs -f object-storage

# Container stats
docker stats object-storage-server

# Health check
curl http://localhost:8080/api/health
```

## 🐛 Troubleshooting

```bash
# Rebuild image
docker-compose build --no-cache

# Restart with fresh start
docker-compose down -v
docker-compose up -d

# Check container logs
docker-compose logs object-storage

# Access container shell for debugging
docker-compose exec object-storage sh
```

## 🌐 URLs

- **Health Check:** http://localhost:8080/api/health
- **Upload:** http://localhost:8080/api/upload
- **Files:** http://localhost:8080/api/files/

---

**For full documentation, see: [DOCKER.md](DOCKER.md)**
