.PHONY: help build up down restart logs clean test

# Default target
help:
	@echo "Object Storage Server - Docker Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build Docker image"
	@echo "  make up             - Start services"
	@echo "  make down           - Stop services"
	@echo "  make restart        - Restart services"
	@echo "  make logs           - View logs"
	@echo "  make logs-f         - Follow logs"
	@echo ""
	@echo "Production:"
	@echo "  make prod-up        - Start production services"
	@echo "  make prod-down      - Stop production services"
	@echo "  make prod-logs      - View production logs"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean          - Remove containers and volumes"
	@echo "  make prune          - Clean Docker system"
	@echo "  make backup         - Backup uploads directory"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run health check"
	@echo "  make shell          - Access container shell"
	@echo ""

# Development commands
build:
	docker-compose build

up:
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "🌐 Server: http://localhost:8080"
	@echo "📊 Health: http://localhost:8080/api/health"

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs

logs-f:
	docker-compose logs -f

# Production commands
prod-up:
	docker-compose -f docker-compose.prod.yml up -d
	@echo "✅ Production services started!"

prod-down:
	docker-compose -f docker-compose.prod.yml down

prod-logs:
	docker-compose -f docker-compose.prod.yml logs -f

prod-restart:
	docker-compose -f docker-compose.prod.yml restart

# Maintenance commands
clean:
	docker-compose down -v
	@echo "✅ Containers and volumes removed"

prune:
	docker system prune -a -f
	@echo "✅ Docker system cleaned"

backup:
	@mkdir -p backups
	tar -czf backups/uploads-$(shell date +%Y%m%d-%H%M%S).tar.gz uploads/
	@echo "✅ Backup created in backups/"

# Testing commands
test:
	@echo "🧪 Testing health endpoint..."
	@curl -s http://localhost:8080/api/health || echo "❌ Server not responding"

shell:
	docker-compose exec object-storage sh

# Stats
stats:
	docker stats --no-stream
