# Command Line Interface Documentation

## Overview

Nutrient Backend API cung cấp một CLI interface mạnh mẽ với nhiều flags và options để config và chạy application.

## Commands

### 1. Server Command

Start the API server với configurable options.

```bash
nutrient-api server [flags]
```

#### Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to configuration file | `--config=./configs/config.prod.yaml` |
| `--env` | `-e` | Environment (dev, staging, prod) | `--env=production` |
| `--port` | `-p` | Server port (overrides config) | `--port=8080` |
| `--host` | | Server host (overrides config) | `--host=0.0.0.0` |
| `--debug` | `-d` | Enable debug mode | `--debug` |
| `--log-level` | `-l` | Log level (debug, info, warn, error) | `--log-level=debug` |
| `--log-format` | | Log format (json, console) | `--log-format=json` |
| `--db-uri` | | MongoDB connection URI | `--db-uri=mongodb://localhost:27017` |
| `--db-name` | | MongoDB database name | `--db-name=nutrient_prod` |
| `--jwt-secret` | | JWT secret key | `--jwt-secret=my-secret-key` |
| `--shutdown-timeout` | | Server shutdown timeout in seconds | `--shutdown-timeout=60` |

#### Examples

```bash
# Start with default config
nutrient-api server

# Start with custom config file
nutrient-api server --config=./configs/config.prod.yaml

# Start with environment override
nutrient-api server --env=production --port=8080

# Start with debug mode
nutrient-api server --debug --log-level=debug

# Start with custom database
nutrient-api server --db-uri=mongodb://localhost:27017 --db-name=nutrient_prod

# Start with custom host and port
nutrient-api server --host=0.0.0.0 --port=3000

# Start with JSON logging
nutrient-api server --log-format=json --log-level=info
```

### 2. Migrate Command

Run database migrations để create indexes và setup collections.

```bash
nutrient-api migrate [flags]
```

#### Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to configuration file | `--config=./configs/config.prod.yaml` |
| `--env` | `-e` | Environment (dev, staging, prod) | `--env=production` |
| `--db-uri` | | MongoDB connection URI | `--db-uri=mongodb://localhost:27017` |
| `--db-name` | | MongoDB database name | `--db-name=nutrient_prod` |
| `--dry-run` | | Show what would be done without executing | `--dry-run` |
| `--force` | `-f` | Skip confirmation prompts | `--force` |

#### Examples

```bash
# Run migrations with default config
nutrient-api migrate

# Run migrations with custom config
nutrient-api migrate --config=./configs/config.prod.yaml

# Run migrations with custom database
nutrient-api migrate --db-uri=mongodb://localhost:27017 --db-name=nutrient_prod

# Dry run (show what would be done without executing)
nutrient-api migrate --dry-run

# Force run (skip confirmation prompts)
nutrient-api migrate --force

# Run migrations in production environment
nutrient-api migrate --env=prod --force
```

### 3. Version Command

Show version information.

```bash
nutrient-api version
```

#### Output

```
Nutrient Backend API
===================
Version:     1.0.0-dev
Go Version:  go1.21.0
OS/Arch:     darwin/arm64
Build Time:  2025-01-15T10:30:00Z
Git Commit:  abc123def456
```

### 4. Info Command

Show detailed application information.

```bash
nutrient-api info
```

#### Output

```
Nutrient Backend API - Application Information
=============================================

Configuration:
  Config Path: ./configs/config.dev.yaml
  Environment: dev

Server:
  Host: localhost
  Port: 8080
  Mode: debug

Database:
  URI:  mongodb://localhost:27017
  Name: nutrient_dev

Features:
  Context Logging:    Enabled
  Response Middleware: Enabled
  JWT Authentication: Enabled
  MongoDB Support:    Enabled
  Excel Import:       Enabled
```

## Global Flags

Các flags này có thể được sử dụng với bất kỳ command nào:

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to configuration file | `--config=./configs/config.yaml` |
| `--env` | `-e` | Environment (dev, staging, prod) | `--env=production` |
| `--debug` | `-d` | Enable debug mode | `--debug` |
| `--log-level` | `-l` | Log level (debug, info, warn, error) | `--log-level=debug` |

## Configuration Priority

Configuration được load theo thứ tự ưu tiên sau:

1. **Command Line Flags** (highest priority)
2. **Environment Variables**
3. **Configuration File**
4. **Default Values** (lowest priority)

### Example Priority

```bash
# This command will:
# 1. Load config from ./configs/config.prod.yaml
# 2. Override with environment variables
# 3. Override with command line flags
nutrient-api server --config=./configs/config.prod.yaml --port=3000 --debug
```

## Environment Variables

Các environment variables có thể được sử dụng để override config:

| Environment Variable | Description | Example |
|---------------------|-------------|---------|
| `APP_ENV` | Environment | `APP_ENV=production` |
| `APP_DEBUG` | Debug mode | `APP_DEBUG=true` |
| `SERVER_HOST` | Server host | `SERVER_HOST=0.0.0.0` |
| `SERVER_PORT` | Server port | `SERVER_PORT=8080` |
| `DB_URI` | Database URI | `DB_URI=mongodb://localhost:27017` |
| `DB_NAME` | Database name | `DB_NAME=nutrient_prod` |
| `JWT_SECRET` | JWT secret | `JWT_SECRET=my-secret-key` |
| `LOG_LEVEL` | Log level | `LOG_LEVEL=debug` |
| `LOG_FORMAT` | Log format | `LOG_FORMAT=json` |

## Docker Usage

### Using Docker Compose

```bash
# Start with default config
docker-compose up

# Start with custom environment
APP_ENV=production docker-compose up

# Start with custom port
SERVER_PORT=3000 docker-compose up
```

### Using Docker Run

```bash
# Run with default config
docker run -p 8080:8080 nutrient-api

# Run with custom config
docker run -p 8080:8080 -v $(pwd)/configs:/app/configs nutrient-api --config=/app/configs/config.prod.yaml

# Run with environment variables
docker run -p 8080:8080 -e APP_ENV=production -e SERVER_PORT=3000 nutrient-api
```

## Development Workflow

### 1. Local Development

```bash
# Start development server
nutrient-api server --debug --log-level=debug

# Run migrations
nutrient-api migrate --env=dev

# Check version
nutrient-api version
```

### 2. Testing

```bash
# Start test server
nutrient-api server --env=test --port=8081 --debug

# Run test migrations
nutrient-api migrate --env=test --force
```

### 3. Production Deployment

```bash
# Start production server
nutrient-api server --env=prod --config=./configs/config.prod.yaml

# Run production migrations
nutrient-api migrate --env=prod --config=./configs/config.prod.yaml --force
```

## Troubleshooting

### Common Issues

1. **Config file not found**
   ```bash
   # Solution: Specify correct config path
   nutrient-api server --config=./configs/config.dev.yaml
   ```

2. **Database connection failed**
   ```bash
   # Solution: Check database URI and ensure MongoDB is running
   nutrient-api server --db-uri=mongodb://localhost:27017
   ```

3. **Port already in use**
   ```bash
   # Solution: Use different port
   nutrient-api server --port=8081
   ```

4. **Permission denied for migrations**
   ```bash
   # Solution: Use force flag or check database permissions
   nutrient-api migrate --force
   ```

### Debug Mode

Enable debug mode để có thêm thông tin:

```bash
# Enable debug mode
nutrient-api server --debug

# Enable debug logging
nutrient-api server --log-level=debug

# Enable debug mode and logging
nutrient-api server --debug --log-level=debug
```

### Logging

Configure logging để debug issues:

```bash
# Console logging (development)
nutrient-api server --log-format=console --log-level=debug

# JSON logging (production)
nutrient-api server --log-format=json --log-level=info

# Error only logging
nutrient-api server --log-level=error
```

## Best Practices

### 1. Use Environment-Specific Configs

```bash
# Development
nutrient-api server --env=dev

# Staging
nutrient-api server --env=staging

# Production
nutrient-api server --env=prod
```

### 2. Use Configuration Files

```bash
# Create environment-specific configs
nutrient-api server --config=./configs/config.dev.yaml
nutrient-api server --config=./configs/config.staging.yaml
nutrient-api server --config=./configs/config.prod.yaml
```

### 3. Use Environment Variables in Production

```bash
# Set environment variables
export APP_ENV=production
export SERVER_PORT=8080
export DB_URI=mongodb://localhost:27017
export DB_NAME=nutrient_prod

# Start server
nutrient-api server
```

### 4. Use Dry Run for Migrations

```bash
# Check what migrations would do
nutrient-api migrate --dry-run

# Run migrations
nutrient-api migrate --force
```

### 5. Use Force Flag for Automated Deployments

```bash
# Automated deployment script
nutrient-api migrate --env=prod --force
nutrient-api server --env=prod
```
