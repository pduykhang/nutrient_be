# API Documentation

## Overview

The Nutrition & Meal Planning API provides comprehensive endpoints for managing nutrition data, meal planning, and shopping lists.

## Base URL

- Development: `http://localhost:8080`
- Production: `https://api.nutrientapp.com`

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message",
  "details": "Additional error details (optional)"
}
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error
- `503` - Service Unavailable

## Endpoints

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "age": 30,
  "weight": 70.0,
  "height": 175.0,
  "gender": "male",
  "goal": "weight_loss"
}
```

**Response:**
```json
{
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "profile": {
      "name": "John Doe",
      "age": 30,
      "weight": 70.0,
      "height": 175.0,
      "gender": "male",
      "goal": "weight_loss"
    },
    "preferences": {
      "language": "en",
      "calorieTarget": 1800.0,
      "macroTargets": {
        "protein": 1.6,
        "carbohydrates": 2.0,
        "fat": 0.8,
        "fiber": 0.03
      }
    }
  },
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": "2025-01-15T10:30:00Z"
}
```

#### Login User
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Food Management

#### Create Food Item
```http
POST /api/v1/foods
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": {
    "en": "Chicken Breast",
    "vi": "Ức Gà"
  },
  "description": {
    "en": "Skinless, boneless chicken breast",
    "vi": "Ức gà không da, không xương"
  },
  "category": "protein",
  "macros": {
    "protein": 31.0,
    "carbohydrates": 0.0,
    "fat": 3.6,
    "fiber": 0.0
  },
  "calories": 165.0,
  "servingSizes": [
    {
      "unit": "gram",
      "amount": 100,
      "gramEquivalent": 100
    },
    {
      "unit": "piece",
      "amount": 1,
      "description": "1 medium breast (174g)",
      "gramEquivalent": 174
    }
  ],
  "visibility": "public"
}
```

#### Search Foods
```http
GET /api/v1/foods/search?q=chicken&lang=vi&limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Food Item
```http
GET /api/v1/foods/{id}
Authorization: Bearer <token>
```

#### Update Food Item
```http
PUT /api/v1/foods/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": {
    "en": "Updated Chicken Breast",
    "vi": "Ức Gà Cập Nhật"
  }
}
```

#### Delete Food Item
```http
DELETE /api/v1/foods/{id}
Authorization: Bearer <token>
```

#### Import Excel
```http
POST /api/v1/foods/import
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <excel-file>
```

### Meal Templates

#### Create Meal Template
```http
POST /api/v1/meal-templates
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My High Protein Breakfast",
  "description": "Perfect breakfast for muscle building",
  "mealType": "breakfast",
  "foodItems": [
    {
      "foodItemId": "507f1f77bcf86cd799439011",
      "servingUnit": "gram",
      "amount": 150
    },
    {
      "foodItemId": "507f1f77bcf86cd799439012",
      "servingUnit": "gram",
      "amount": 100
    }
  ],
  "tags": ["high-protein", "balanced"],
  "isPublic": false
}
```

#### List Meal Templates
```http
GET /api/v1/meal-templates?mealType=breakfast&limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Meal Template
```http
GET /api/v1/meal-templates/{id}
Authorization: Bearer <token>
```

#### Update Meal Template
```http
PUT /api/v1/meal-templates/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Breakfast Template"
}
```

#### Delete Meal Template
```http
DELETE /api/v1/meal-templates/{id}
Authorization: Bearer <token>
```

### Meal Plans

#### Create Meal Plan
```http
POST /api/v1/meal-plans
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Weight Loss Week 1",
  "description": "First week of my weight loss journey",
  "startDate": "2025-01-06T00:00:00Z",
  "endDate": "2025-01-12T23:59:59Z",
  "planType": "weekly",
  "goal": "weight_loss",
  "targetCalories": 1800,
  "targetMacros": {
    "protein": 135,
    "carbohydrates": 180,
    "fat": 60,
    "fiber": 30
  }
}
```

#### List Meal Plans
```http
GET /api/v1/meal-plans?planType=weekly&limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Meal Plan
```http
GET /api/v1/meal-plans/{id}
Authorization: Bearer <token>
```

#### Update Meal Plan
```http
PUT /api/v1/meal-plans/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Meal Plan",
  "targetCalories": 1900
}
```

#### Delete Meal Plan
```http
DELETE /api/v1/meal-plans/{id}
Authorization: Bearer <token>
```

### Shopping Lists

#### Generate Shopping List
```http
POST /api/v1/shopping-lists/generate/{mealPlanId}
Authorization: Bearer <token>
```

#### List Shopping Lists
```http
GET /api/v1/shopping-lists?limit=10&offset=0
Authorization: Bearer <token>
```

#### Toggle Shopping Item
```http
PUT /api/v1/shopping-lists/{id}/items/{itemId}/check
Authorization: Bearer <token>
Content-Type: application/json

{
  "checked": true
}
```

### Reports

#### Weekly Report
```http
GET /api/v1/reports/weekly?startDate=2025-01-06&endDate=2025-01-12
Authorization: Bearer <token>
```

#### Monthly Report
```http
GET /api/v1/reports/monthly?month=2025-01&year=2025
Authorization: Bearer <token>
```

### Health Checks

#### Liveness Probe
```http
GET /health/liveness
```

**Response:**
```json
{
  "status": "UP",
  "timestamp": 1705123456,
  "service": "nutrient-api"
}
```

#### Readiness Probe
```http
GET /health/readiness
```

**Response:**
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

## Data Models

### Food Item
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": {
    "en": "Chicken Breast",
    "vi": "Ức Gà"
  },
  "searchTerms": ["chicken breast", "uc ga", "ức gà"],
  "description": {
    "en": "Skinless, boneless chicken breast",
    "vi": "Ức gà không da, không xương"
  },
  "category": "protein",
  "macros": {
    "protein": 31.0,
    "carbohydrates": 0.0,
    "fat": 3.6,
    "fiber": 0.0,
    "sugar": 0.0
  },
  "micros": {
    "vitaminA": 0.0,
    "vitaminC": 0.0,
    "calcium": 15.0,
    "iron": 0.7,
    "sodium": 74.0,
    "potassium": 256.0
  },
  "servingSizes": [
    {
      "unit": "gram",
      "amount": 100,
      "description": "100 grams",
      "gramEquivalent": 100
    }
  ],
  "calories": 165.0,
  "createdBy": "507f191e810c19729de860ea",
  "visibility": "public",
  "source": "user",
  "imageUrl": "https://example.com/chicken-breast.jpg",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

### Meal Template
```json
{
  "id": "507f1f77bcf86cd799439022",
  "userId": "507f191e810c19729de860ea",
  "name": "My High Protein Breakfast",
  "description": "Perfect breakfast for muscle building",
  "mealType": "breakfast",
  "foodItems": [
    {
      "foodItemId": "507f1f77bcf86cd799439011",
      "foodName": "Chicken Breast",
      "servingUnit": "gram",
      "amount": 150,
      "calories": 247.5,
      "macros": {
        "protein": 46.5,
        "carbohydrates": 0.0,
        "fat": 5.4,
        "fiber": 0.0
      }
    }
  ],
  "totalCalories": 247.5,
  "totalMacros": {
    "protein": 46.5,
    "carbohydrates": 0.0,
    "fat": 5.4,
    "fiber": 0.0
  },
  "tags": ["high-protein", "balanced"],
  "isPublic": false,
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

### Meal Plan
```json
{
  "id": "507f1f77bcf86cd799439033",
  "userId": "507f191e810c19729de860ea",
  "name": "Weight Loss Week 1",
  "description": "First week of my weight loss journey",
  "startDate": "2025-01-06T00:00:00Z",
  "endDate": "2025-01-12T23:59:59Z",
  "planType": "weekly",
  "goal": "weight_loss",
  "targetCalories": 1800,
  "targetMacros": {
    "protein": 135,
    "carbohydrates": 180,
    "fat": 60,
    "fiber": 30
  },
  "dailyMeals": [
    {
      "date": "2025-01-06T00:00:00Z",
      "dayOfWeek": "Monday",
      "meals": [
        {
          "id": "meal_001",
          "mealType": "breakfast",
          "time": "07:00",
          "templateId": "507f1f77bcf86cd799439022",
          "foodItems": [
            {
              "foodItemId": "507f1f77bcf86cd799439011",
              "foodName": "Chicken Breast",
              "foodCategory": "protein",
              "servingUnit": "gram",
              "amount": 150,
              "calories": 247.5,
              "macros": {
                "protein": 46.5,
                "carbohydrates": 0.0,
                "fat": 5.4,
                "fiber": 0.0
              }
            }
          ],
          "calories": 247.5,
          "macros": {
            "protein": 46.5,
            "carbohydrates": 0.0,
            "fat": 5.4,
            "fiber": 0.0
          },
          "notes": "",
          "isCompleted": false
        }
      ],
      "totalCalories": 247.5,
      "totalMacros": {
        "protein": 46.5,
        "carbohydrates": 0.0,
        "fat": 5.4,
        "fiber": 0.0
      },
      "notes": "First day, feeling excited!",
      "isCompleted": false
    }
  ],
  "totalCalories": 1732.5,
  "status": "active",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

### Shopping List
```json
{
  "id": "507f1f77bcf86cd799439044",
  "userId": "507f191e810c19729de860ea",
  "mealPlanId": "507f1f77bcf86cd799439033",
  "items": [
    {
      "foodItemId": "507f1f77bcf86cd799439011",
      "foodName": "Chicken Breast",
      "totalAmount": 1050,
      "unit": "gram",
      "checked": false
    },
    {
      "foodItemId": "507f1f77bcf86cd799439012",
      "foodName": "Brown Rice",
      "totalAmount": 700,
      "unit": "gram",
      "checked": true
    }
  ],
  "totalCost": 0,
  "status": "pending",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

## Rate Limiting

- **Authentication endpoints**: 10 requests per minute per IP
- **Other endpoints**: 100 requests per minute per user
- **Bulk operations**: 5 requests per minute per user

## Pagination

List endpoints support pagination with these query parameters:

- `limit`: Number of items per page (default: 20, max: 100)
- `offset`: Number of items to skip (default: 0)

**Example:**
```http
GET /api/v1/foods/search?q=chicken&limit=10&offset=20
```

## Filtering and Sorting

### Food Items
- `category`: Filter by food category (protein, vegetable, fruit, dairy, grain)
- `visibility`: Filter by visibility (public, private)
- `source`: Filter by source (user, imported)

### Meal Templates
- `mealType`: Filter by meal type (breakfast, lunch, dinner, snack)
- `isPublic`: Filter by public/private templates

### Meal Plans
- `planType`: Filter by plan type (weekly, monthly)
- `status`: Filter by status (draft, active, completed)
- `goal`: Filter by goal (weight_loss, muscle_gain, maintenance)

## Search

### Food Search
The food search endpoint supports multi-language search with the following features:

- **Text search**: Searches in food names and search terms
- **Language support**: Supports English and Vietnamese
- **Fuzzy matching**: Handles typos and variations
- **Category filtering**: Filter results by food category

**Search Query Parameters:**
- `q`: Search query (required)
- `lang`: Language preference (en, vi)
- `category`: Filter by category
- `limit`: Number of results (default: 20)
- `offset`: Pagination offset (default: 0)

**Example:**
```http
GET /api/v1/foods/search?q=chicken&lang=vi&category=protein&limit=10
```

## Webhooks (Future Feature)

The API will support webhooks for real-time notifications:

- `meal_plan.created`
- `meal_plan.updated`
- `meal_plan.completed`
- `shopping_list.generated`
- `user.goal_achieved`

## SDKs and Libraries

### JavaScript/TypeScript
```bash
npm install @nutrient-app/api-client
```

### Python
```bash
pip install nutrient-api-client
```

### Go
```bash
go get github.com/nutrient-app/go-client
```

## Support

For API support and questions:

- **Documentation**: https://docs.nutrientapp.com
- **Email**: api-support@nutrientapp.com
- **GitHub Issues**: https://github.com/nutrient-app/api/issues
