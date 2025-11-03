# User Authentication & Management Architecture

## Overview

Hệ thống đã được tách thành 2 phần chính:
1. **User Authentication (Auth)**: Xử lý đăng ký, đăng nhập, logout, refresh token, validate token
2. **User Management (User)**: Quản lý thông tin user bao gồm profile, preferences, và targets

## Architecture Design

```
┌─────────────────────────────────────────────────────────────┐
│                      API Layer                               │
│  ┌──────────────────────┐  ┌────────────────────────────┐  │
│  │   AuthHandler        │  │    UserHandler              │  │
│  │                      │  │                            │  │
│  │ - Register           │  │ - GetProfile                │  │
│  │ - Login              │  │ - UpdateProfile              │  │
│  │ - Logout              │  │ - UpdatePreferences          │  │
│  │ - Refresh             │  │ - ChangePassword             │  │
│  │ - Validate            │  │                            │  │
│  └──────────────────────┘  └────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                             │
│  ┌──────────────────────┐  ┌────────────────────────────┐  │
│  │   AuthService        │  │    UserService             │  │
│  │                      │  │                            │  │
│  │ - Register           │  │ - GetProfile                │  │
│  │ - Login              │  │ - UpdateProfile              │  │
│  │ - RefreshToken       │  │ - UpdatePreferences          │  │
│  │ - ValidateToken      │  │ - ChangePassword             │  │
│  │ - generateTokens     │  │ - calculateCalorieTarget    │  │
│  └──────────────────────┘  │ - calculateMacroTargets     │  │
│                            └────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Repository Layer                            │
│                    UserRepository                            │
│                                                            │
│ - Create                                                   │
│ - GetByID                                                  │
│ - GetByEmail                                               │
│ - Update                                                   │
│ - Delete                                                   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                            │
│                        User Entity                           │
│                                                            │
│ - ID                                                       │
│ - Email                                                    │
│ - PasswordHash                                             │
│ - Profile (Name, Age, Weight, Height, Gender, Goal)       │
│ - Preferences (Language, CalorieTarget, MacroTargets)     │
│ - CreatedAt, UpdatedAt                                     │
└─────────────────────────────────────────────────────────────┘
```

## Component Relationships

### AuthService
- **Responsibility**: Authentication logic only
  - User registration (email + password only)
  - User login
  - Token generation (access + refresh)
  - Token validation
  - Token refresh
- **Dependencies**: UserRepository (to create/get users for auth)

### UserService
- **Responsibility**: User data management
  - Get user profile
  - Update user profile
  - Update user preferences
  - Change password
  - Calculate calorie targets based on profile
  - Calculate macro targets based on goal
- **Dependencies**: UserRepository

### Shared Components
- **UserRepository**: Used by both AuthService and UserService
- **Domain.User**: Shared entity model
- **DTOs**: Separate request/response DTOs for each service

## API Endpoints

### Public Auth Endpoints (No Authentication Required)
```
POST /api/v1/auth/register
  Body: { "email": string, "password": string }
  
POST /api/v1/auth/login
  Body: { "email": string, "password": string }
  
POST /api/v1/auth/refresh
  Body: { "refreshToken": string }
  
POST /api/v1/auth/validate
  Headers: { "Authorization": "Bearer <token>" }
```

### Protected User Endpoints (Authentication Required)
```
GET /api/v1/users/profile
  Returns: User profile and preferences
  
PUT /api/v1/users/profile
  Body: { "name"?: string, "age"?: int, "weight"?: float, 
          "height"?: float, "gender"?: string, "goal"?: string }
  
PUT /api/v1/users/preferences
  Body: { "language"?: string, "calorieTarget"?: float,
          "macroTargets"?: MacroNutrients }
  
PUT /api/v1/users/password
  Body: { "currentPassword": string, "newPassword": string }
  
POST /api/v1/auth/logout
```

## Registration Flow

### Step 1: Register (Auth)
```
User → POST /auth/register { email, password }
  → AuthService.Register()
  → Creates user with empty profile
  → Returns tokens
```

### Step 2: Complete Profile (User)
```
User → PUT /users/profile { name, age, weight, height, gender, goal }
  → UserService.UpdateProfile()
  → Calculates calorie target and macro targets
  → Updates user profile
```

## Design Decisions

### Why Separate Auth and User Services?

1. **Separation of Concerns**
   - Authentication logic is separate from user data management
   - Clear boundaries between responsibilities

2. **Scalability**
   - Can scale authentication and user management independently
   - Easier to add features like OAuth, 2FA to AuthService
   - User management features don't affect auth logic

3. **Security**
   - Auth operations are isolated
   - User profile updates require authentication but are separate from auth flow

4. **Flexibility**
   - Can register with minimal info (email/password)
   - User can complete profile later
   - Profile updates don't require re-authentication

### Registration Design

**Before**: Registration required all profile information upfront
```json
{
  "email": "...",
  "password": "...",
  "name": "...",
  "age": ...,
  "weight": ...,
  "height": ...,
  "gender": "...",
  "goal": "..."
}
```

**After**: Two-step process
1. Register with email/password only
2. Complete profile separately via user endpoints

This allows:
- Faster registration
- User can complete profile at their convenience
- Profile can be updated anytime without affecting authentication

### Target Calculation

Calorie and macro targets are automatically calculated when:
- User sets/updates their goal
- User updates weight/height/age (and goal is already set)

The calculation logic is in `UserService` because it's part of user data management, not authentication.

## DTOs Structure

### Request DTOs
- `request.RegisterRequest`: Email + password only
- `request.LoginRequest`: Email + password
- `request.UpdateProfileRequest`: Optional profile fields
- `request.UpdatePreferencesRequest`: Optional preference fields
- `request.ChangePasswordRequest`: Current + new password

### Response DTOs
- `response.UserResponse`: Complete user data (without password)
- `response.UserProfileResponse`: Profile information
- `response.UserPreferencesResponse`: Preferences and targets

## Error Handling

- **AuthService errors**: Authentication-related (invalid credentials, token errors)
- **UserService errors**: User data errors (user not found, validation errors)
- Both services return domain-specific errors

## Future Enhancements

1. **Token Blacklist**: For logout functionality
2. **Password Reset**: Can be added to AuthService
3. **Email Verification**: Can be added to AuthService
4. **Profile Picture**: Can be added to UserService
5. **Activity Level**: Can affect calorie calculations

