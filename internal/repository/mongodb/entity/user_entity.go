package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
)

// UserEntity represents a user in MongoDB
type UserEntity struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
	Profile      UserProfileEntity  `bson:"profile"`
	Preferences  UserPreferencesEntity `bson:"preferences"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

// UserProfileEntity represents user profile in MongoDB
type UserProfileEntity struct {
	Name   string  `bson:"name"`
	Age    int     `bson:"age"`
	Weight float64 `bson:"weight"`
	Height float64 `bson:"height"`
	Gender string  `bson:"gender"`
	Goal   string  `bson:"goal"`
}

// UserPreferencesEntity represents user preferences in MongoDB
type UserPreferencesEntity struct {
	Language      string              `bson:"language"`
	CalorieTarget float64             `bson:"calorieTarget"`
	MacroTargets  MacroNutrientsEntity `bson:"macroTargets"`
}

// MacroNutrientsEntity represents macronutrient values in MongoDB
type MacroNutrientsEntity struct {
	Protein       float64 `bson:"protein"`
	Carbohydrates float64 `bson:"carbohydrates"`
	Fat           float64 `bson:"fat"`
	Fiber         float64 `bson:"fiber"`
	Sugar         float64 `bson:"sugar,omitempty"`
}

// ToDomain converts UserEntity to domain.User
func (e *UserEntity) ToDomain() *domain.User {
	if e == nil {
		return nil
	}

	return &domain.User{
		ID:           e.ID,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		Profile: domain.UserProfile{
			Name:   e.Profile.Name,
			Age:    e.Profile.Age,
			Weight: e.Profile.Weight,
			Height: e.Profile.Height,
			Gender: e.Profile.Gender,
			Goal:   e.Profile.Goal,
		},
		Preferences: domain.UserPreferences{
			Language:      e.Preferences.Language,
			CalorieTarget: e.Preferences.CalorieTarget,
			MacroTargets: domain.MacroNutrients{
				Protein:       e.Preferences.MacroTargets.Protein,
				Carbohydrates: e.Preferences.MacroTargets.Carbohydrates,
				Fat:           e.Preferences.MacroTargets.Fat,
				Fiber:         e.Preferences.MacroTargets.Fiber,
				Sugar:         e.Preferences.MacroTargets.Sugar,
			},
		},
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// FromDomain creates UserEntity from domain.User
func (e *UserEntity) FromDomain(u *domain.User) error {
	if u == nil {
		return nil
	}

	// Use domain ID directly (still primitive.ObjectID for now)
	id := u.ID
	if u.ID.IsZero() {
		id = primitive.NewObjectID()
	}

	e.ID = id
	e.Email = u.Email
	e.PasswordHash = u.PasswordHash
	e.Profile = UserProfileEntity{
		Name:   u.Profile.Name,
		Age:    u.Profile.Age,
		Weight: u.Profile.Weight,
		Height: u.Profile.Height,
		Gender: u.Profile.Gender,
		Goal:   u.Profile.Goal,
	}
	e.Preferences = UserPreferencesEntity{
		Language:      u.Preferences.Language,
		CalorieTarget: u.Preferences.CalorieTarget,
		MacroTargets: MacroNutrientsEntity{
			Protein:       u.Preferences.MacroTargets.Protein,
			Carbohydrates: u.Preferences.MacroTargets.Carbohydrates,
			Fat:           u.Preferences.MacroTargets.Fat,
			Fiber:         u.Preferences.MacroTargets.Fiber,
			Sugar:         u.Preferences.MacroTargets.Sugar,
		},
	}
	e.CreatedAt = u.CreatedAt
	e.UpdatedAt = u.UpdatedAt

	return nil
}

