# Nutrition & Meal Planning Backend

A comprehensive Golang backend system for nutrition tracking and meal planning with MongoDB storage, multi-language food search, meal plan generation, shopping list support, and nutritional reporting capabilities.

## Features

- **User Authentication**: JWT-based authentication with registration and login
- **Food Management**: CRUD operations for food items with multi-language support
- **Meal Templates**: Reusable meal combinations for quick meal planning
- **Meal Plans**: Weekly/monthly eating schedules with calorie and macro calculations
- **Shopping Lists**: Auto-generated shopping lists from meal plans
- **Nutrition Reports**: Weekly/monthly nutrition statistics and goal tracking
- **Excel Import**: Import food data from Excel files
- **Health Checks**: Kubernetes-ready liveness and readiness probes
- **Docker Support**: Complete Docker setup for local development

## Technology Stack

- **Framework**: Gin (HTTP router)
- **CLI**: Cobra (command-line interface)
- **Configuration**: Viper (configuration management)
- **Database**: MongoDB with official driver
- **Logging**: Zap with custom abstraction layer
- **Authentication**: JWT tokens
- **Validation**: go-playground/validator
- **Excel**: excelize for import/export
- **Message Queue**: NATS (for future async operations)

## Project Structure

```
nutrient_be/
├── cmd/api/                    # Application entry points
│   ├── main.go                # Root command
│   ├── server.go              # Server command
│   └── migrate.go             # Migration command
├── internal/
│   ├── domain/                # Core business entities
│   │   ├── user.go           # User domain
│   │   ├── meal.go           # Meal and meal plan domains
│   │   └── shopping.go       # Shopping list domain
│   ├── repository/           # Data access layer
│   │   ├── interfaces.go     # Repository interfaces
│   │   └── mongodb/          # MongoDB implementations
│   ├── service/              # Business logic layer
│   │   ├── auth_service.go   # Authentication service
│   │   └── placeholder_services.go # Other services
│   ├── handler/              # HTTP handlers
│   │   ├── rest/             # REST API handlers
│   │   └── middleware/       # HTTP middleware
│   ├── pkg/                  # Internal packages
│   │   └── logger/          # Logging abstraction
│   ├── config/               # Configuration management
│   └── database/             # Database connection
├── configs/                  # Configuration files
├── deployments/              # Docker and deployment files
├── scripts/                  # Utility scripts
└── docs/                     # Documentation
```

## Quick Start

### Prerequisites

- Go 1.24+
- MongoDB 7.0+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd nutrient_be
   go mod download
   ```

2. **Start MongoDB**:
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:7.0
   
   # Or using local MongoDB installation
   mongod
   ```

3. **Run migrations**:
   ```bash
   go run cmd/api/main.go migrate --config=configs/config.dev.yaml
   ```

4. **Start the server**:
   ```bash
   go run cmd/api/main.go server --config=configs/config.dev.yaml
   ```

5. **Test the API**:
   ```bash
   # Health check
   curl http://localhost:8080/health/liveness
   
   # Register a user
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "test@example.com",
       "password": "password123",
       "name": "Test User",
       "age": 30,
       "weight": 70,
       "height": 175,
       "gender": "male",
       "goal": "weight_loss"
     }'
   ```

### Docker Development

1. **Start all services**:
   ```bash
   make docker-up
   ```

2. **Run migrations**:
   ```bash
   make migrate
   ```

3. **Access services**:
   - API: http://localhost:8080
   - MongoDB Express: http://localhost:8081 (admin/admin123)
   - NATS Monitoring: http://localhost:8222

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Food Management
- `POST /api/v1/foods` - Create food item
- `GET /api/v1/foods/search?q=query&lang=vi` - Search foods
- `GET /api/v1/foods/:id` - Get food item
- `PUT /api/v1/foods/:id` - Update food item
- `DELETE /api/v1/foods/:id` - Delete food item
- `POST /api/v1/foods/import` - Import from Excel

### Meal Templates
- `POST /api/v1/meal-templates` - Create template
- `GET /api/v1/meal-templates` - List templates
- `GET /api/v1/meal-templates/:id` - Get template
- `PUT /api/v1/meal-templates/:id` - Update template
- `DELETE /api/v1/meal-templates/:id` - Delete template

### Meal Plans
- `POST /api/v1/meal-plans` - Create meal plan
- `GET /api/v1/meal-plans` - List meal plans
- `GET /api/v1/meal-plans/:id` - Get meal plan
- `PUT /api/v1/meal-plans/:id` - Update meal plan
- `DELETE /api/v1/meal-plans/:id` - Delete meal plan

### Shopping Lists
- `POST /api/v1/shopping-lists/generate/:mealPlanId` - Generate from meal plan
- `GET /api/v1/shopping-lists` - List shopping lists
- `PUT /api/v1/shopping-lists/:id/items/:itemId/check` - Toggle item checked

### Reports
- `GET /api/v1/reports/weekly?startDate=2025-01-01` - Weekly nutrition report
- `GET /api/v1/reports/monthly?month=2025-01` - Monthly nutrition report

### Health Checks
- `GET /health/liveness` - Liveness probe
- `GET /health/readiness` - Readiness probe

## Configuration

The application uses Viper for configuration management with support for YAML files and environment variables.

### Configuration Structure

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release
  read_timeout: 10
  write_timeout: 10
  shutdown_timeout: 30

database:
  uri: "mongodb://localhost:27017"
  database: "nutrient_db"
  max_pool_size: 100
  min_pool_size: 10
  connect_timeout: 10

auth:
  jwt_secret: "${JWT_SECRET}"
  jwt_expiration: 3600
  refresh_expiration: 604800

logger:
  level: "debug"  # debug, info, warn, error
  development: true
  encoding: "console"  # console, json
```

### Environment Variables

- `JWT_SECRET` - Secret key for JWT tokens
- `MONGODB_URI` - MongoDB connection string
- `NATS_URL` - NATS server URL
- `CONFIG_PATH` - Path to configuration file

## Database Schema

### Collections

- **users**: User accounts and profiles
- **foods**: Food items with nutrition data
- **meal_templates**: Reusable meal combinations
- **meal_plans**: Complete meal schedules
- **shopping_lists**: Generated shopping lists

### Indexes

- Text search on food names and search terms
- User-based queries for all collections
- Date-based queries for meal plans
- Category-based queries for foods

## Logging

The application uses a custom logging abstraction that wraps Zap, allowing easy switching between different logging libraries without code changes.

```go
// Usage example
logger.Info("User registered", 
    logger.String("email", user.Email),
    logger.String("userID", user.ID.Hex()))
```

## Development Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Start Docker services
make docker-up

# Stop Docker services
make docker-down

# Run database migrations
make migrate

# Run linters
make lint
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test` and `make lint`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [ ] Complete food service implementation
- [ ] Complete meal template service
- [ ] Complete meal plan generation
- [ ] Complete shopping list generation
- [ ] Complete reporting service
- [ ] Add comprehensive tests
- [ ] Add API documentation with Swagger
- [ ] Add rate limiting
- [ ] Add caching layer
- [ ] Add metrics and monitoring
