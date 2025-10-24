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

// mealTemplateRepository implements repository.MealTemplateRepository
type mealTemplateRepository struct {
	collection *mongo.Collection
}

// NewMealTemplateRepository creates a new meal template repository
func NewMealTemplateRepository(db *mongo.Database) repository.MealTemplateRepository {
	return &mealTemplateRepository{
		collection: db.Collection("meal_templates"),
	}
}

// Create creates a new meal template
func (r *mealTemplateRepository) Create(ctx context.Context, template *domain.MealTemplate) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("failed to create meal template: %w", err)
	}
	return nil
}

// GetByID retrieves a meal template by ID
func (r *mealTemplateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealTemplate, error) {
	var template domain.MealTemplate
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("meal template not found")
		}
		return nil, fmt.Errorf("failed to get meal template: %w", err)
	}
	return &template, nil
}

// GetByUser retrieves meal templates by user
func (r *mealTemplateRepository) GetByUser(ctx context.Context, userID primitive.ObjectID, mealType string, limit, offset int) ([]*domain.MealTemplate, error) {
	filter := bson.M{"userId": userID}
	if mealType != "" {
		filter["mealType"] = mealType
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal templates: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*domain.MealTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("failed to decode meal templates: %w", err)
	}

	return templates, nil
}

// GetPublicTemplates retrieves public meal templates
func (r *mealTemplateRepository) GetPublicTemplates(ctx context.Context, mealType string, limit, offset int) ([]*domain.MealTemplate, error) {
	filter := bson.M{"isPublic": true}
	if mealType != "" {
		filter["mealType"] = mealType
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get public meal templates: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*domain.MealTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("failed to decode meal templates: %w", err)
	}

	return templates, nil
}

// Update updates a meal template
func (r *mealTemplateRepository) Update(ctx context.Context, template *domain.MealTemplate) error {
	template.UpdatedAt = time.Now()

	filter := bson.M{"_id": template.ID}
	update := bson.M{"$set": template}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update meal template: %w", err)
	}
	return nil
}

// Delete deletes a meal template
func (r *mealTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete meal template: %w", err)
	}
	return nil
}
