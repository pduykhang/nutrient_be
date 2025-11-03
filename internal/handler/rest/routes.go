package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nutrient_be/internal/handler/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, handlers *Handlers) {
	// Add global middleware first
	r.Use(middleware.LoggingMiddleware(handlers.Auth.logger))
	r.Use(middleware.RecoveryMiddleware(handlers.Auth.logger))
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ContextMiddleware(handlers.Auth.logger))  // Add context middleware
	r.Use(middleware.ResponseMiddleware(handlers.Auth.logger)) // Add response middleware

	// Health checks (no auth required)
	r.HEAD("/health/liveness", handlers.Health.Liveness)
	r.GET("/health/readiness", handlers.Health.Readiness)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Auth.Register)
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/refresh", handlers.Auth.Refresh)
			auth.POST("/validate", handlers.Auth.Validate)
		}

		// Protected routes (auth required)
		protected := v1.Group("")
		protected.Use(middleware.DefaultUserAuthMiddleware(handlers.Auth.logger))
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("/profile", handlers.User.GetProfile)
				users.PUT("/profile", handlers.User.UpdateProfile)
				users.PUT("/preferences", handlers.User.UpdatePreferences)
				users.PUT("/password", handlers.User.ChangePassword)
			}

			// Auth (protected)
			authProtected := protected.Group("/auth")
			{
				authProtected.POST("/logout", handlers.Auth.Logout)
			}

			// Foods
			foods := protected.Group("/foods")
			{
				foods.POST("", handlers.Food.Create)
				foods.GET("/search", handlers.Food.Search)
				foods.GET("/:id", handlers.Food.Get)
				foods.PUT("/:id", handlers.Food.Update)
				foods.DELETE("/:id", handlers.Food.Delete)
				foods.POST("/import", handlers.Food.ImportExcel)
			}

			// Meal templates
			templates := protected.Group("/meal-templates")
			{
				templates.POST("", handlers.Meal.CreateTemplate)
				templates.GET("", handlers.Meal.ListTemplates)
				templates.GET("/:id", handlers.Meal.GetTemplate)
				templates.PUT("/:id", handlers.Meal.UpdateTemplate)
				templates.DELETE("/:id", handlers.Meal.DeleteTemplate)
			}

			// Meal plans
			plans := protected.Group("/meal-plans")
			{
				plans.POST("", handlers.MealPlan.Create)
				plans.GET("", handlers.MealPlan.List)
				plans.GET("/:id", handlers.MealPlan.Get)
				plans.PUT("/:id", handlers.MealPlan.Update)
				plans.DELETE("/:id", handlers.MealPlan.Delete)
			}

			// Shopping lists
			shopping := protected.Group("/shopping-lists")
			{
				shopping.POST("/generate/:mealPlanId", handlers.Shopping.Generate)
				shopping.GET("", handlers.Shopping.List)
				shopping.PUT("/:id/items/:itemId/check", handlers.Shopping.ToggleItem)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/weekly", handlers.Report.Weekly)
				reports.GET("/monthly", handlers.Report.Monthly)
			}
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})
}
