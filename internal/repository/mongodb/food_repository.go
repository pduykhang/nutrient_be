package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"nutrient_be/internal/domain"
)

// foodRepository handles food data operations
type foodRepository struct {
	collection *mongo.Collection
}

// NewFoodRepository creates a new food repository
func NewFoodRepository(db *mongo.Database) *foodRepository {
	return &foodRepository{
		collection: db.Collection("foods"),
	}
}

// Create creates a new food item
func (r *foodRepository) Create(ctx context.Context, food *domain.FoodItem) error {
	food.CreatedAt = time.Now()
	food.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, food)
	if err != nil {
		return fmt.Errorf("failed to create food item: %w", err)
	}
	return nil
}

// GetByID retrieves a food item by ID
func (r *foodRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.FoodItem, error) {
	var food domain.FoodItem
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&food)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("food item not found")
		}
		return nil, fmt.Errorf("failed to get food item: %w", err)
	}
	return &food, nil
}

// Search searches for food items using text search
func (r *foodRepository) Search(ctx context.Context, query string, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error) {
	// Normalize search query
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	// Build search filter
	filter := bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"searchTerms": bson.M{"$regex": normalizedQuery, "$options": "i"}},
					{"name.en": bson.M{"$regex": normalizedQuery, "$options": "i"}},
					{"name.vi": bson.M{"$regex": normalizedQuery, "$options": "i"}},
				},
			},
			{
				"$or": []bson.M{
					{"visibility": "public"},
					{"createdBy": userID},
				},
			},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search food items: %w", err)
	}
	defer cursor.Close(ctx)

	var foods []*domain.FoodItem
	if err := cursor.All(ctx, &foods); err != nil {
		return nil, fmt.Errorf("failed to decode food items: %w", err)
	}

	return foods, nil
}

// GetByCategory retrieves food items by category
func (r *foodRepository) GetByCategory(ctx context.Context, category string, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error) {
	filter := bson.M{
		"category": category,
		"$or": []bson.M{
			{"visibility": "public"},
			{"createdBy": userID},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get food items by category: %w", err)
	}
	defer cursor.Close(ctx)

	var foods []*domain.FoodItem
	if err := cursor.All(ctx, &foods); err != nil {
		return nil, fmt.Errorf("failed to decode food items: %w", err)
	}

	return foods, nil
}

// GetByUser retrieves food items created by a specific user
func (r *foodRepository) GetByUser(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error) {
	filter := bson.M{"createdBy": userID}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get user food items: %w", err)
	}
	defer cursor.Close(ctx)

	var foods []*domain.FoodItem
	if err := cursor.All(ctx, &foods); err != nil {
		return nil, fmt.Errorf("failed to decode food items: %w", err)
	}

	return foods, nil
}

// Update updates a food item
func (r *foodRepository) Update(ctx context.Context, food *domain.FoodItem) error {
	food.UpdatedAt = time.Now()

	filter := bson.M{"_id": food.ID}
	update := bson.M{"$set": food}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update food item: %w", err)
	}
	return nil
}

// Delete deletes a food item
func (r *foodRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete food item: %w", err)
	}
	return nil
}

// GetPublicFoods retrieves public food items
func (r *foodRepository) GetPublicFoods(ctx context.Context, limit, offset int) ([]*domain.FoodItem, error) {
	filter := bson.M{"visibility": "public"}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get public food items: %w", err)
	}
	defer cursor.Close(ctx)

	var foods []*domain.FoodItem
	if err := cursor.All(ctx, &foods); err != nil {
		return nil, fmt.Errorf("failed to decode food items: %w", err)
	}

	return foods, nil
}
