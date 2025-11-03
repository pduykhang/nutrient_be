# Docker Development Setup Guide

Hướng dẫn chi tiết để chạy ứng dụng với Docker Compose và hot reload (Air).

## Prerequisites

- Docker Engine >= 20.10
- Docker Compose >= 2.0
- Make (optional, để dùng các commands ngắn gọn)

## Quick Start

### 1. Start tất cả services với hot reload

```bash
make docker-dev
```

Hoặc:

```bash
cd deployments
docker-compose -f docker-compose.yml up --build
```

### 2. View logs

```bash
make docker-logs
```

Hoặc:

```bash
docker-compose -f deployments/docker-compose.yml logs -f api
```

### 3. Stop services

```bash
make docker-down
```

Hoặc:

```bash
docker-compose -f deployments/docker-compose.yml down
```

## Services

### 1. API Service (nutrient-api-dev)

- **Port**: `8080`
- **Hot Reload**: Enabled (Air)
- **Config**: `configs/config.dev.yaml`
- **Health Check**: `http://localhost:8080/health/liveness`

### 2. MongoDB (mongo)

- **Port**: `27017`
- **Database**: `nutrient_db`
- **Init Script**: `scripts/init-mongo.js`

### 3. Mongo Express (mongo-express)

- **Port**: `8081`
- **URL**: `http://localhost:8081`
- **Username**: `admin`
- **Password**: `admin123`

### 4. NATS (nats)

- **Client Port**: `4222`
- **HTTP Port**: `8222`
- **Monitoring**: `http://localhost:8222`

## Configuration

### Config File: `configs/config.dev.yaml`

```yaml
server:
  host: "0.0.0.0"  # Listen on all interfaces
  port: 8080
  mode: "debug"

database:
  uri: "mongodb://mongo:27017"  # Docker service name
  database: "nutrient_db"

auth:
  jwt_secret: "dev-secret-key"
  jwt_expiration: 3600

nats:
  url: "nats://nats:4222"  # Docker service name
  enabled: true

logger:
  level: "debug"
  development: true
  encoding: "console"
```

### Environment Variables Override

Viper hỗ trợ override config từ environment variables:

- `DATABASE_URI` -> `database.uri`
- `AUTH_JWT_SECRET` -> `auth.jwt_secret`
- `SERVER_PORT` -> `server.port`
- `NATS_URL` -> `nats.url`
- `NATS_ENABLED` -> `nats.enabled`

Format: `SECTION_KEY` maps to `section.key` trong config.

Ví dụ:
```bash
# Override database URI
export DATABASE_URI=mongodb://custom-host:27017

# Override JWT secret
export AUTH_JWT_SECRET=my-custom-secret
```

## Hot Reload với Air

### Cách hoạt động

1. Edit bất kỳ `.go` file nào trong project
2. Air tự động phát hiện changes
3. Air rebuild application
4. Application tự động restart
5. Ready cho changes tiếp theo

### Air Configuration: `.air.toml`

```toml
[build]
  cmd = "go build -o ./tmp/nutrient-api ./cmd/api/..."
  bin = "./tmp/nutrient-api"
  args_bin = ["server", "--config=configs/config.dev.yaml"]
  
  exclude_dir = ["tmp", "bin", "vendor", ".git"]
  include_ext = [".go"]
```

### Troubleshooting Hot Reload

**Problem**: Changes không trigger rebuild

**Solutions**:
1. Check file extension có trong `include_ext`
2. Check directory có trong `exclude_dir`
3. View Air logs: `docker exec -it nutrient-api-dev cat tmp/air.log`
4. Check volume mount: `docker inspect nutrient-api-dev | grep Mounts`

## Volume Mapping

```yaml
volumes:
  - ..:/app              # Entire project mounted
  - /app/tmp             # Excluded (Air creates here)
  - /app/bin             # Excluded
```

- Source code được mount từ host vào `/app` trong container
- `tmp/` và `bin/` được exclude để tránh conflicts
- Changes trên host sẽ sync vào container ngay lập tức

## Network

Tất cả services trong network `nutrient-network`:

```yaml
networks:
  nutrient-network:
    driver: bridge
```

Services có thể communicate bằng service names:
- `mongodb://mongo:27017`
- `nats://nats:4222`

## Health Checks

### Liveness Probe

```bash
curl http://localhost:8080/health/liveness
```

Response:
```json
{
  "status": "UP",
  "timestamp": 1705123456,
  "service": "nutrient-api"
}
```

### Readiness Probe

```bash
curl http://localhost:8080/health/readiness
```

Response:
```json
{
  "status": "UP",
  "timestamp": 1705123456,
  "checks": {
    "database": "UP"
  },
  "service": "nutrient-api"
}
```

## Common Commands

### Start services in background

```bash
make docker-up
# hoặc
docker-compose -f deployments/docker-compose.yml up -d
```

### View all logs

```bash
docker-compose -f deployments/docker-compose.yml logs -f
```

### Rebuild API service

```bash
make docker-rebuild
# hoặc
docker-compose -f deployments/docker-compose.yml build --no-cache api
```

### Execute commands in container

```bash
# Enter container
docker exec -it nutrient-api-dev sh

# Run migrate
docker exec -it nutrient-api-dev ./tmp/nutrient-api migrate --config=configs/config.dev.yaml
```

### Check service status

```bash
docker-compose -f deployments/docker-compose.yml ps
```

## Database Access

### MongoDB

```bash
# Connect from host
mongosh mongodb://localhost:27017/nutrient_db

# Connect from container
docker exec -it nutrient-api-dev mongosh mongodb://mongo:27017/nutrient_db
```

### Mongo Express (Web UI)

1. Open browser: `http://localhost:8081`
2. Login với username: `admin`, password: `admin123`

## API Endpoints

Base URL: `http://localhost:8080`

### Health
- `GET /health/liveness`
- `GET /health/readiness`

### Authentication
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`

## Development Workflow

1. **Start services**:
   ```bash
   make docker-dev
   ```

2. **Edit code** trong bất kỳ `.go` file nào

3. **Watch logs** để xem hot reload:
   ```bash
   make docker-logs
   ```

4. **Test API**:
   ```bash
   curl http://localhost:8080/health/liveness
   ```

5. **Stop services**:
   ```bash
   make docker-down
   ```

## Troubleshooting

### Container không start

```bash
# Check logs
docker-compose -f deployments/docker-compose.yml logs api

# Check container status
docker ps -a | grep nutrient
```

### Database connection failed

1. Check MongoDB service đang chạy:
   ```bash
   docker-compose -f deployments/docker-compose.yml ps mongo
   ```

2. Check MongoDB logs:
   ```bash
   docker-compose -f deployments/docker-compose.yml logs mongo
   ```

3. Verify connection string trong config:
   ```yaml
   database:
     uri: "mongodb://mongo:27017"  # Must use service name
   ```

### Port already in use

```bash
# Check what's using port 8080
lsof -i :8080

# Kill process or change port in docker-compose.yml
```

### Hot reload không hoạt động

1. Check Air logs:
   ```bash
   docker exec -it nutrient-api-dev cat tmp/air.log
   ```

2. Verify volume mount:
   ```bash
   docker inspect nutrient-api-dev | grep Mounts -A 10
   ```

3. Check file permissions

4. Rebuild container:
   ```bash
   make docker-rebuild
   ```

## Production vs Development

### Development (current setup)

- Hot reload enabled (Air)
- Debug mode
- Console logging
- Volume mounts for live changes

### Production

Use `Dockerfile` (not `Dockerfile.dev`):
- No hot reload
- Release mode
- JSON logging
- Optimized binary

## Clean Up

### Remove containers and volumes

```bash
docker-compose -f deployments/docker-compose.yml down -v
```

### Remove images

```bash
docker-compose -f deployments/docker-compose.yml down --rmi all
```

## Resources

- [Air Documentation](https://github.com/air-verse/air)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Viper Configuration](https://github.com/spf13/viper)
