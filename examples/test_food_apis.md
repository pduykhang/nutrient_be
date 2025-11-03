# Food API Testing Guide

This guide provides examples for testing the Food APIs: Create, Search, and Get by ID.

## Prerequisites

1. Start the API server (default: `http://localhost:8080`)
2. Have a registered user account (or register first)

## Step 1: Authenticate and Get Token

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "..."
}
```

Save the `accessToken` for subsequent requests.

## Step 2: Create a Food Item

```bash
curl -X POST http://localhost:8080/api/v1/foods \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": {
      "en": "Grilled Chicken Breast",
      "vi": "Ức gà nướng"
    },
    "searchTerms": ["chicken", "breast", "grilled", "protein"],
    "description": {
      "en": "Lean protein source, perfect for muscle building",
      "vi": "Nguồn protein nạc, hoàn hảo cho việc xây dựng cơ bắp"
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
      "vitaminA": 20.0,
      "vitaminC": 0.0,
      "calcium": 15.0,
      "iron": 0.7,
      "sodium": 84.0,
      "potassium": 256.0
    },
    "servingSizes": [
      {
        "unit": "gram",
        "amount": 100,
        "description": "100 grams",
        "gramEquivalent": 100
      },
      {
        "unit": "piece",
        "amount": 1,
        "description": "1 piece (about 200g)",
        "gramEquivalent": 200
      }
    ],
    "calories": 165.0,
    "visibility": "public",
    "imageUrl": "https://example.com/chicken-breast.jpg"
  }'
```

**Response:**
```json
{
  "message": "Food created successfully",
  "data": {
    "message": "Food created successfully"
  }
}
```

### Example: Create a Vegetable

```bash
curl -X POST http://localhost:8080/api/v1/foods \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": {
      "en": "Broccoli",
      "vi": "Bông cải xanh"
    },
    "searchTerms": ["broccoli", "vegetable", "green"],
    "category": "vegetable",
    "macros": {
      "protein": 2.8,
      "carbohydrates": 7.0,
      "fat": 0.4,
      "fiber": 2.6,
      "sugar": 1.5
    },
    "micros": {
      "vitaminA": 31.0,
      "vitaminC": 89.2,
      "calcium": 47.0,
      "iron": 0.7,
      "sodium": 33.0,
      "potassium": 316.0
    },
    "servingSizes": [
      {
        "unit": "gram",
        "amount": 100,
        "description": "100 grams",
        "gramEquivalent": 100
      },
      {
        "unit": "cup",
        "amount": 1,
        "description": "1 cup chopped (91g)",
        "gramEquivalent": 91
      }
    ],
    "calories": 34.0,
    "visibility": "public"
  }'
```

## Step 3: Search Food Items

```bash
curl -X GET "http://localhost:8080/api/v1/foods/search?query=chicken&limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Query Parameters:**
- `query` (required): Search term
- `limit` (optional, default: 20): Maximum number of results
- `offset` (optional, default: 0): Number of results to skip

**Response:**
```json
{
  "data": [
    {
      "id": "65a1b2c3d4e5f6789012345",
      "name": {
        "en": "Grilled Chicken Breast",
        "vi": "Ức gà nướng"
      },
      "searchTerms": ["chicken", "breast", "grilled", "protein"],
      "category": "protein",
      "macros": {
        "protein": 31.0,
        "carbohydrates": 0.0,
        "fat": 3.6,
        "fiber": 0.0,
        "sugar": 0.0
      },
      "servingSizes": [...],
      "calories": 165.0,
      "createdBy": "65a1b2c3d4e5f6789012346",
      "visibility": "public",
      "source": "user",
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Food search successful"
}
```

### More Search Examples

```bash
# Search for protein foods
curl -X GET "http://localhost:8080/api/v1/foods/search?query=protein&limit=5" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Search for vegetables
curl -X GET "http://localhost:8080/api/v1/foods/search?query=vegetable&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Search with pagination
curl -X GET "http://localhost:8080/api/v1/foods/search?query=chicken&limit=5&offset=5" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Step 4: Get Food by ID

```bash
# Replace FOOD_ID with actual food ID from search results
curl -X GET http://localhost:8080/api/v1/foods/FOOD_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "data": {
    "id": "65a1b2c3d4e5f6789012345",
    "name": {
      "en": "Grilled Chicken Breast",
      "vi": "Ức gà nướng"
    },
    "searchTerms": ["chicken", "breast", "grilled", "protein"],
    "description": {
      "en": "Lean protein source, perfect for muscle building",
      "vi": "Nguồn protein nạc, hoàn hảo cho việc xây dựng cơ bắp"
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
      "vitaminA": 20.0,
      "vitaminC": 0.0,
      "calcium": 15.0,
      "iron": 0.7,
      "sodium": 84.0,
      "potassium": 256.0
    },
    "servingSizes": [
      {
        "unit": "gram",
        "amount": 100,
        "description": "100 grams",
        "gramEquivalent": 100
      },
      {
        "unit": "piece",
        "amount": 1,
        "description": "1 piece (about 200g)",
        "gramEquivalent": 200
      }
    ],
    "calories": 165.0,
    "createdBy": "65a1b2c3d4e5f6789012346",
    "visibility": "public",
    "source": "user",
    "imageUrl": "https://example.com/chicken-breast.jpg",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  },
  "message": "Food retrieved successfully"
}
```

**Error Response (404):**
```json
{
  "error": "Food item not found",
  "message": "Food item not found"
}
```

## Using the Test Script

You can also use the automated test script:

```bash
# Set environment variables if needed
export BASE_URL=http://localhost:8080
export TEST_EMAIL=your-email@example.com
export TEST_PASSWORD=your-password

# Run the test script
./examples/test_food_apis.sh
```

## Valid Categories

- `protein`
- `vegetable`
- `fruit`
- `dairy`
- `grain`

## Valid Visibility

- `public`: Food is visible to all users
- `private`: Food is only visible to the creator

## Notes

1. All food endpoints require authentication (JWT token)
2. Search returns both public foods and foods created by the authenticated user
3. Food IDs are MongoDB ObjectIDs (24-character hex strings)
4. Multi-language support: Currently supports `en` (English) and `vi` (Vietnamese)
5. Serving sizes must have at least one entry
6. Macros and micros values are per 100g by default

