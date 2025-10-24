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

// shoppingListRepository implements repository.ShoppingListRepository
type shoppingListRepository struct {
	collection *mongo.Collection
}

// NewShoppingListRepository creates a new shopping list repository
func NewShoppingListRepository(db *mongo.Database) repository.ShoppingListRepository {
	return &shoppingListRepository{
		collection: db.Collection("shopping_lists"),
	}
}

// Create creates a new shopping list
func (r *shoppingListRepository) Create(ctx context.Context, list *domain.ShoppingList) error {
	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, list)
	if err != nil {
		return fmt.Errorf("failed to create shopping list: %w", err)
	}
	return nil
}

// GetByID retrieves a shopping list by ID
func (r *shoppingListRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.ShoppingList, error) {
	var list domain.ShoppingList
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("shopping list not found")
		}
		return nil, fmt.Errorf("failed to get shopping list: %w", err)
	}
	return &list, nil
}

// GetByUser retrieves shopping lists by user
func (r *shoppingListRepository) GetByUser(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*domain.ShoppingList, error) {
	filter := bson.M{"userId": userID}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get shopping lists: %w", err)
	}
	defer cursor.Close(ctx)

	var lists []*domain.ShoppingList
	if err := cursor.All(ctx, &lists); err != nil {
		return nil, fmt.Errorf("failed to decode shopping lists: %w", err)
	}

	return lists, nil
}

// GetByMealPlan retrieves a shopping list by meal plan ID
func (r *shoppingListRepository) GetByMealPlan(ctx context.Context, mealPlanID primitive.ObjectID) (*domain.ShoppingList, error) {
	var list domain.ShoppingList
	err := r.collection.FindOne(ctx, bson.M{"mealPlanId": mealPlanID}).Decode(&list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("shopping list not found")
		}
		return nil, fmt.Errorf("failed to get shopping list by meal plan: %w", err)
	}
	return &list, nil
}

// Update updates a shopping list
func (r *shoppingListRepository) Update(ctx context.Context, list *domain.ShoppingList) error {
	list.UpdatedAt = time.Now()

	filter := bson.M{"_id": list.ID}
	update := bson.M{"$set": list}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update shopping list: %w", err)
	}
	return nil
}

// Delete deletes a shopping list
func (r *shoppingListRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete shopping list: %w", err)
	}
	return nil
}

// ToggleItemChecked toggles the checked status of a shopping list item
func (r *shoppingListRepository) ToggleItemChecked(ctx context.Context, listID primitive.ObjectID, itemID primitive.ObjectID, checked bool) error {
	filter := bson.M{
		"_id":              listID,
		"items.foodItemId": itemID,
	}

	update := bson.M{
		"$set": bson.M{
			"items.$.checked": checked,
			"updatedAt":       time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to toggle item checked status: %w", err)
	}
	return nil
}
