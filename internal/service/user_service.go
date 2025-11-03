package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/dto/response"
	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user profile and preferences management
type UserService struct {
	userRepo UserRepository
	logger   logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo UserRepository, log logger.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   log,
	}
}

// GetProfile retrieves user profile by ID
func (s *UserService) GetProfile(ctx context.Context, userID string) (*response.UserResponse, error) {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userIDObj)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	s.logger.Info(ctx, "User profile retrieved", logger.String("userID", userID))
	return domainUserToResponse(user), nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *request.UpdateProfileRequest) (*response.UserResponse, error) {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userIDObj)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update profile fields
	if req.Name != nil {
		user.Profile.Name = *req.Name
	}
	if req.Age != nil {
		user.Profile.Age = *req.Age
	}
	if req.Weight != nil {
		user.Profile.Weight = *req.Weight
	}
	if req.Height != nil {
		user.Profile.Height = *req.Height
	}
	if req.Gender != nil {
		user.Profile.Gender = *req.Gender
	}
	if req.Goal != nil {
		user.Profile.Goal = *req.Goal
		// Recalculate calorie target and macro targets when goal changes
		if user.Profile.Weight > 0 && user.Profile.Height > 0 && user.Profile.Age > 0 {
			user.Preferences.CalorieTarget = calculateCalorieTarget(
				user.Profile.Weight,
				user.Profile.Height,
				user.Profile.Age,
				user.Profile.Gender,
				user.Profile.Goal,
			)
			user.Preferences.MacroTargets = calculateMacroTargets(user.Profile.Goal)
		}
	}

	// Also recalculate if weight/height/age changes and goal is set
	if (req.Weight != nil || req.Height != nil || req.Age != nil) && user.Profile.Goal != "" {
		if user.Profile.Weight > 0 && user.Profile.Height > 0 && user.Profile.Age > 0 {
			user.Preferences.CalorieTarget = calculateCalorieTarget(
				user.Profile.Weight,
				user.Profile.Height,
				user.Profile.Age,
				user.Profile.Gender,
				user.Profile.Goal,
			)
		}
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	s.logger.Info(ctx, "User profile updated", logger.String("userID", userID))
	return domainUserToResponse(user), nil
}

// UpdatePreferences updates user preferences
func (s *UserService) UpdatePreferences(ctx context.Context, userID string, req *request.UpdatePreferencesRequest) (*response.UserResponse, error) {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userIDObj)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update preferences fields
	if req.Language != nil {
		user.Preferences.Language = *req.Language
	}
	if req.CalorieTarget != nil {
		user.Preferences.CalorieTarget = *req.CalorieTarget
	}
	if req.MacroTargets != nil {
		user.Preferences.MacroTargets = domain.MacroNutrients{
			Protein:       req.MacroTargets.Protein,
			Carbohydrates: req.MacroTargets.Carbohydrates,
			Fat:           req.MacroTargets.Fat,
			Fiber:         req.MacroTargets.Fiber,
			Sugar:         req.MacroTargets.Sugar,
		}
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user preferences: %w", err)
	}

	s.logger.Info(ctx, "User preferences updated", logger.String("userID", userID))
	return domainUserToResponse(user), nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordRequest) error {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userIDObj)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Info(ctx, "User password changed", logger.String("userID", userID))
	return nil
}

// calculateCalorieTarget calculates daily calorie target based on user profile
func calculateCalorieTarget(weight, height float64, age int, gender, goal string) float64 {
	// Basic BMR calculation (Mifflin-St Jeor Equation)
	var bmr float64
	if gender == "male" {
		bmr = 10*weight + 6.25*height - 5*float64(age) + 5
	} else {
		bmr = 10*weight + 6.25*height - 5*float64(age) - 161
	}

	// Activity factor (sedentary)
	activityFactor := 1.2
	maintenanceCalories := bmr * activityFactor

	// Adjust based on goal
	switch goal {
	case "weight_loss":
		return maintenanceCalories - 500 // 500 calorie deficit
	case "muscle_gain":
		return maintenanceCalories + 300 // 300 calorie surplus
	default: // maintenance
		return maintenanceCalories
	}
}

// calculateMacroTargets calculates macro targets based on goal
func calculateMacroTargets(goal string) domain.MacroNutrients {
	switch goal {
	case "weight_loss":
		return domain.MacroNutrients{
			Protein:       1.6, // g per kg body weight
			Carbohydrates: 2.0,
			Fat:           0.8,
			Fiber:         0.03,
		}
	case "muscle_gain":
		return domain.MacroNutrients{
			Protein:       2.2,
			Carbohydrates: 4.0,
			Fat:           1.0,
			Fiber:         0.03,
		}
	default: // maintenance
		return domain.MacroNutrients{
			Protein:       1.8,
			Carbohydrates: 3.0,
			Fat:           0.9,
			Fiber:         0.03,
		}
	}
}

