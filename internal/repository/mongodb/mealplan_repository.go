package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/repository"
)

// mealPlanRepository implements repository.MealPlanRepository
type mealPlanRepository struct {
	collection *mongo.Collection
}

// NewMealPlanRepository creates a new meal plan repository
func NewMealPlanRepository(db *mongo.Database) repository.MealPlanRepository {
	return &mealPlanRepository{
		collection: db.Collection("meal_plans"),
	}
}

// Create creates a new meal plan
func (r *mealPlanRepository) Create(ctx context.Context, plan *domain.MealPlan) error {
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, plan)
	if err != nil {
		return fmt.Errorf("failed to create meal plan: %w", err)
	}
	return nil
}

// GetByID retrieves a meal plan by ID
func (r *mealPlanRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealPlan, error) {
	var plan domain.MealPlan
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&plan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("meal plan not found")
		}
		return nil, fmt.Errorf("failed to get meal plan: %w", err)
	}
	return &plan, nil
}

// GetByUser retrieves meal plans by user
func (r *mealPlanRepository) GetByUser(ctx context.Context, userID primitive.ObjectID, planType string, limit, offset int) ([]*domain.MealPlan, error) {
	filter := bson.M{"userId": userID}
	if planType != "" {
		filter["planType"] = planType
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"startDate": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal plans: %w", err)
	}
	defer cursor.Close(ctx)

	var plans []*domain.MealPlan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, fmt.Errorf("failed to decode meal plans: %w", err)
	}

	return plans, nil
}

// GetByUserAndDateRange retrieves meal plans by user and date range
func (r *mealPlanRepository) GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate string) ([]*domain.MealPlan, error) {
	filter := bson.M{
		"userId": userID,
		"$or": []bson.M{
			{
				"startDate": bson.M{"$lte": endDate},
				"endDate":   bson.M{"$gte": startDate},
			},
		},
	}

	opts := options.Find().SetSort(bson.M{"startDate": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal plans by date range: %w", err)
	}
	defer cursor.Close(ctx)

	var plans []*domain.MealPlan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, fmt.Errorf("failed to decode meal plans: %w", err)
	}

	return plans, nil
}

// Update updates a meal plan
func (r *mealPlanRepository) Update(ctx context.Context, plan *domain.MealPlan) error {
	plan.UpdatedAt = time.Now()

	filter := bson.M{"_id": plan.ID}
	update := bson.M{"$set": plan}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update meal plan: %w", err)
	}
	return nil
}

// Delete deletes a meal plan
func (r *mealPlanRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete meal plan: %w", err)
	}
	return nil
}

// UpdateMealCompletion updates meal completion status
func (r *mealPlanRepository) UpdateMealCompletion(ctx context.Context, planID primitive.ObjectID, mealID string, isCompleted bool) error {
	filter := bson.M{
		"_id":                 planID,
		"dailyMeals.meals.id": mealID,
	}

	update := bson.M{
		"$set": bson.M{
			"dailyMeals.$[].meals.$[meal].isCompleted": isCompleted,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"meal.id": mealID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update meal completion: %w", err)
	}
	return nil
}
