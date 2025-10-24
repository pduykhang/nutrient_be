# Architecture Documentation

## System Overview

The Nutrition & Meal Planning Backend is built using clean architecture principles with clear separation of concerns and dependency inversion. The system is designed to be scalable, maintainable, and testable.

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)

The core business entities and rules. This layer is independent of external frameworks and contains:

- **User**: User accounts and profiles
- **FoodItem**: Food ingredients with nutrition data
- **MealTemplate**: Reusable meal combinations
- **MealPlan**: Complete eating schedules
- **ShoppingList**: Generated shopping lists

**Key Principles:**
- No external dependencies
- Pure business logic
- Rich domain models with behavior
- Value objects for complex data types

### 2. Repository Layer (`internal/repository/`)

Data access abstraction that defines interfaces for data operations:

- **UserRepository**: User data operations
- **FoodRepository**: Food item operations with search
- **MealTemplateRepository**: Meal template operations
- **MealPlanRepository**: Meal plan operations
- **ShoppingListRepository**: Shopping list operations

**Key Principles:**
- Interface-based design
- Database-agnostic interfaces
- Clear separation of concerns
- Easy testing with mocks

### 3. Service Layer (`internal/service/`)

Business logic orchestration and application services:

- **AuthService**: Authentication and authorization
- **FoodService**: Food management and search
- **MealService**: Meal template management
- **MealPlanService**: Meal plan generation and management
- **ShoppingService**: Shopping list generation
- **ReportService**: Nutrition reporting

**Key Principles:**
- Orchestrates domain operations
- Handles business workflows
- Manages transactions
- Coordinates between repositories

### 4. Handler Layer (`internal/handler/`)

HTTP request/response handling and API endpoints:

- **REST Handlers**: HTTP API endpoints
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)
- **Validation**: Request validation and sanitization
- **Serialization**: JSON marshaling/unmarshaling

**Key Principles:**
- Thin layer over services
- HTTP-specific concerns only
- Input validation and sanitization
- Error handling and status codes

### 5. Infrastructure Layer

External dependencies and infrastructure concerns:

- **Database**: MongoDB connection and configuration
- **Logger**: Structured logging with Zap
- **Config**: Configuration management with Viper
- **CLI**: Command-line interface with Cobra

## Data Flow

```
HTTP Request → Handler → Service → Repository → Database
     ↓           ↓        ↓         ↓           ↓
HTTP Response ← Handler ← Service ← Repository ← Database
```

### Request Processing Flow

1. **HTTP Request** arrives at Gin router
2. **Middleware** processes cross-cutting concerns (auth, logging)
3. **Handler** validates input and calls appropriate service
4. **Service** orchestrates business logic and calls repositories
5. **Repository** performs data operations on MongoDB
6. **Response** flows back through the layers

## Database Design

### MongoDB Collections

#### Users Collection
```javascript
{
  _id: ObjectId,
  email: String (unique),
  passwordHash: String,
  profile: {
    name: String,
    age: Number,
    weight: Number,
    height: Number,
    gender: String,
    goal: String
  },
  preferences: {
    language: String,
    calorieTarget: Number,
    macroTargets: Object
  },
  createdAt: Date,
  updatedAt: Date
}
```

#### Foods Collection
```javascript
{
  _id: ObjectId,
  name: Object, // Multi-language: {en: "Chicken", vi: "Gà"}
  searchTerms: [String], // Normalized search terms
  description: Object, // Multi-language descriptions
  category: String, // protein, vegetable, fruit, dairy, grain
  macros: {
    protein: Number,
    carbohydrates: Number,
    fat: Number,
    fiber: Number,
    sugar: Number
  },
  micros: {
    vitaminA: Number,
    vitaminC: Number,
    calcium: Number,
    iron: Number,
    sodium: Number,
    potassium: Number
  },
  servingSizes: [{
    unit: String,
    amount: Number,
    description: String,
    gramEquivalent: Number
  }],
  calories: Number, // Per 100g
  createdBy: ObjectId,
  visibility: String, // public, private
  source: String, // user, imported
  imageUrl: String,
  createdAt: Date,
  updatedAt: Date
}
```

#### Meal Templates Collection
```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  name: String,
  description: String,
  mealType: String, // breakfast, lunch, dinner, snack
  foodItems: [{
    foodItemId: ObjectId,
    foodName: String, // Denormalized
    servingUnit: String,
    amount: Number,
    calories: Number, // Calculated
    macros: Object // Calculated
  }],
  totalCalories: Number, // Calculated
  totalMacros: Object, // Calculated
  tags: [String],
  isPublic: Boolean,
  createdAt: Date,
  updatedAt: Date
}
```

#### Meal Plans Collection
```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  name: String,
  description: String,
  startDate: Date,
  endDate: Date,
  planType: String, // weekly, monthly
  goal: String, // weight_loss, muscle_gain, maintenance
  targetCalories: Number,
  targetMacros: Object,
  dailyMeals: [{
    date: Date,
    dayOfWeek: String,
    meals: [{
      id: String,
      mealType: String,
      time: String,
      templateId: ObjectId,
      foodItems: [{
        foodItemId: ObjectId,
        foodName: String, // Denormalized
        foodCategory: String, // Denormalized
        servingUnit: String,
        amount: Number,
        calories: Number, // Calculated
        macros: Object // Calculated
      }],
      calories: Number, // Calculated
      macros: Object, // Calculated
      notes: String,
      isCompleted: Boolean
    }],
    totalCalories: Number, // Calculated
    totalMacros: Object, // Calculated
    notes: String,
    isCompleted: Boolean
  }],
  totalCalories: Number, // Calculated
  status: String, // draft, active, completed
  createdAt: Date,
  updatedAt: Date
}
```

#### Shopping Lists Collection
```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  mealPlanId: ObjectId,
  items: [{
    foodItemId: ObjectId,
    foodName: String,
    totalAmount: Number,
    unit: String,
    checked: Boolean
  }],
  totalCost: Number,
  status: String, // pending, completed
  createdAt: Date,
  updatedAt: Date
}
```

### Indexes

#### Performance Indexes
```javascript
// Users
db.users.createIndex({ "email": 1 }, { unique: true })

// Foods
db.foods.createIndex({ "searchTerms": "text" })
db.foods.createIndex({ "createdBy": 1, "visibility": 1 })
db.foods.createIndex({ "category": 1 })
db.foods.createIndex({ "source": 1 })

// Meal Templates
db.meal_templates.createIndex({ "userId": 1, "mealType": 1 })
db.meal_templates.createIndex({ "userId": 1, "isPublic": 1 })

// Meal Plans
db.meal_plans.createIndex({ "userId": 1, "startDate": -1 })
db.meal_plans.createIndex({ "userId": 1, "planType": 1 })
db.meal_plans.createIndex({ "userId": 1, "status": 1 })

// Shopping Lists
db.shopping_lists.createIndex({ "userId": 1, "mealPlanId": 1 })
db.shopping_lists.createIndex({ "userId": 1, "status": 1 })
```

## Security Architecture

### Authentication & Authorization

1. **JWT Tokens**: Stateless authentication with access and refresh tokens
2. **Password Hashing**: bcrypt with salt for secure password storage
3. **Middleware**: Request-level authentication and authorization
4. **Rate Limiting**: Protection against abuse and DoS attacks

### Data Security

1. **Input Validation**: All inputs validated and sanitized
2. **SQL Injection Prevention**: MongoDB parameterized queries
3. **CORS**: Cross-origin resource sharing configuration
4. **Environment Variables**: Sensitive data in environment variables

## Error Handling

### Error Types

1. **Validation Errors**: Input validation failures (400)
2. **Authentication Errors**: Invalid credentials or tokens (401)
3. **Authorization Errors**: Insufficient permissions (403)
4. **Not Found Errors**: Resource not found (404)
5. **Conflict Errors**: Resource conflicts (409)
6. **Server Errors**: Internal server errors (500)
7. **Service Unavailable**: External service failures (503)

### Error Response Format

```json
{
  "error": "Error message",
  "details": "Additional error details",
  "code": "ERROR_CODE",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

## Logging Strategy

### Structured Logging

Using Zap logger with structured fields:

```go
logger.Info("User registered successfully",
    logger.String("email", user.Email),
    logger.String("userID", user.ID.Hex()),
    logger.String("ip", clientIP))
```

### Log Levels

- **DEBUG**: Detailed debugging information
- **INFO**: General information about application flow
- **WARN**: Warning messages for potential issues
- **ERROR**: Error messages for failed operations
- **FATAL**: Critical errors that cause application shutdown

### Log Fields

- **Request ID**: Unique identifier for each request
- **User ID**: User identifier for authenticated requests
- **IP Address**: Client IP address
- **User Agent**: Client user agent
- **Duration**: Request processing time
- **Status Code**: HTTP response status code

## Configuration Management

### Configuration Sources

1. **Default Values**: Hardcoded defaults in code
2. **YAML Files**: Configuration files for different environments
3. **Environment Variables**: Runtime configuration overrides
4. **Command Line**: CLI parameter overrides

### Configuration Structure

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"
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
  level: "debug"
  development: true
  encoding: "console"
```

## Testing Strategy

### Test Types

1. **Unit Tests**: Individual function and method testing
2. **Integration Tests**: Service layer integration testing
3. **Repository Tests**: Database operation testing
4. **Handler Tests**: HTTP endpoint testing
5. **End-to-End Tests**: Complete workflow testing

### Test Structure

```
tests/
├── unit/           # Unit tests
├── integration/    # Integration tests
├── e2e/           # End-to-end tests
├── fixtures/      # Test data fixtures
└── mocks/         # Mock implementations
```

## Deployment Architecture

### Container Strategy

1. **Multi-stage Dockerfile**: Optimized build and runtime images
2. **Docker Compose**: Local development environment
3. **Health Checks**: Kubernetes-ready health probes
4. **Graceful Shutdown**: Proper application shutdown handling

### Environment Configuration

1. **Development**: Local development with hot reload
2. **Staging**: Pre-production testing environment
3. **Production**: Production environment with optimizations

## Monitoring & Observability

### Metrics

1. **Request Metrics**: Request count, duration, status codes
2. **Business Metrics**: User registrations, meal plans created
3. **System Metrics**: CPU, memory, database connections
4. **Error Metrics**: Error rates and types

### Health Checks

1. **Liveness Probe**: Application is running
2. **Readiness Probe**: Application is ready to serve requests
3. **Database Health**: Database connectivity check
4. **External Services**: External service availability

## Scalability Considerations

### Horizontal Scaling

1. **Stateless Design**: No server-side session storage
2. **Load Balancing**: Multiple application instances
3. **Database Sharding**: MongoDB sharding for large datasets
4. **Caching**: Redis for frequently accessed data

### Performance Optimization

1. **Database Indexes**: Optimized query performance
2. **Connection Pooling**: Efficient database connections
3. **Denormalization**: Reduced join operations
4. **Async Processing**: Background task processing

## Future Enhancements

### Planned Features

1. **Real-time Notifications**: WebSocket support
2. **Caching Layer**: Redis integration
3. **Message Queue**: NATS for async processing
4. **Microservices**: Service decomposition
5. **API Gateway**: Centralized API management
6. **GraphQL**: Alternative API interface
7. **Mobile SDKs**: Native mobile app support

### Technical Debt

1. **Complete Service Implementation**: Finish remaining services
2. **Comprehensive Testing**: Increase test coverage
3. **API Documentation**: OpenAPI/Swagger documentation
4. **Rate Limiting**: Implement rate limiting middleware
5. **Metrics Collection**: Prometheus metrics integration
6. **Tracing**: Distributed tracing with Jaeger
