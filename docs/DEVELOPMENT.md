# Unified IT Operations Portal — Development Guide

## Local Development Setup

### Prerequisites
- Docker & Docker Compose v2
- Go 1.22+
- Node.js 20+ (for web frontend)

### Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/AlexCortada/UnifiedMC.git
cd UnifiedMC

# 2. Start data infrastructure
docker compose -f infrastructure/docker/docker-compose.yml up -d

# 3. Wait for services to be ready
sleep 15

# 4. Run database migrations
for f in database/migrations/*.sql; do
  PGPASSWORD=*** psql -h localhost -U unifiedmc -d unifiedmc -f "$f"
done

# 5. Verify infrastructure
curl http://localhost:5432  # PostgreSQL
curl http://localhost:6379  # Redis
curl http://localhost:9200  # Elasticsearch

# 6. Start API gateway
cd services/api-gateway && go run .
```

### Infrastructure Services

| Service | URL | Credentials |
|---|---|---|
| PostgreSQL | localhost:5432 | unifiedmc / unifiedmc_dev |
| Redis | localhost:6379 | (none) |
| Kafka | localhost:9092 | (none) |
| Elasticsearch | localhost:9200 | (none) |
| ClickHouse | localhost:8123 | (none) |

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| SERVICE_NAME | uop-api-gateway | Service identifier |
| PORT | 8080 | API gateway port |
| DATABASE_URL | postgresql://... | PostgreSQL connection |
| REDIS_URL | redis://localhost:6379 | Redis connection |
| KAFKA_BROKERS | localhost:9092 | Kafka brokers |
| ELASTICSEARCH_URL | http://localhost:9200 | Elasticsearch URL |

### Development Commands

```bash
# Run all Go tests
go test ./...

# Run a specific service
cd services/api-gateway && go run .

# Build all services
go build ./services/... ./connectors/...

# Format code
gofmt -w ./...
```
