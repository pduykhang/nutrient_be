package rest

import (
	"go.mongodb.org/mongo-driver/mongo"

	"nutrient_be/internal/config"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	Auth     *AuthHandler
	User     *UserHandler
	Health   *HealthHandler
	Food     *FoodHandler
	Meal     *MealHandler
	MealPlan *MealPlanHandler
	Shopping *ShoppingHandler
	Report   *ReportHandler
}

// NewHandlers creates a new handlers instance
func NewHandlers(
	authService *service.AuthService,
	userService *service.UserService,
	foodService *service.FoodService,
	mealService *service.MealService,
	mealPlanService *service.MealPlanService,
	shoppingService *service.ShoppingService,
	reportService *service.ReportService,
	db *mongo.Client,
	log logger.Logger,
	cfg config.Config,
) *Handlers {
	return &Handlers{
		Auth:     NewAuthHandler(authService, log, cfg.Auth),
		User:     NewUserHandler(userService, log),
		Health:   NewHealthHandler(db, log),
		Food:     NewFoodHandler(foodService, log),
		Meal:     NewMealHandler(mealService, log),
		MealPlan: NewMealPlanHandler(mealPlanService, log),
		Shopping: NewShoppingHandler(shoppingService, log),
		Report:   NewReportHandler(reportService, log),
	}
}
