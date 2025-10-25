package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"nutrient_be/internal/config"
	"nutrient_be/internal/database"
	"nutrient_be/internal/handler/rest"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/repository/mongodb"
	"nutrient_be/internal/service"
)

// Global flags
var (
	configPath      string
	environment     string
	port            int
	host            string
	debug           bool
	logLevel        string
	logFormat       string
	dbURI           string
	dbName          string
	jwtSecret       string
	shutdownTimeout int
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API server",
	Long: `Start the nutrient backend API server with configurable options.

Examples:
  # Start with default config
  nutrient-api server

  # Start with custom config file
  nutrient-api server --config=./configs/config.prod.yaml

  # Start with environment override
  nutrient-api server --env=production --port=8080

  # Start with debug mode
  nutrient-api server --debug --log-level=debug

  # Start with custom database
  nutrient-api server --db-uri=mongodb://localhost:27017 --db-name=nutrient_prod`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Configuration flags
	serverCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file (e.g., ./configs/config.yaml)")
	serverCmd.Flags().StringVarP(&environment, "env", "e", "", "Environment (dev, staging, prod)")
	serverCmd.Flags().IntVarP(&port, "port", "p", 0, "Server port (overrides config)")
	serverCmd.Flags().StringVarP(&host, "host", "", "", "Server host (overrides config)")
	serverCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	serverCmd.Flags().StringVarP(&logLevel, "log-level", "l", "", "Log level (debug, info, warn, error)")
	serverCmd.Flags().StringVarP(&logFormat, "log-format", "", "", "Log format (json, console)")
	serverCmd.Flags().StringVarP(&dbURI, "db-uri", "", "", "MongoDB connection URI (overrides config)")
	serverCmd.Flags().StringVarP(&dbName, "db-name", "", "", "MongoDB database name (overrides config)")
	serverCmd.Flags().StringVarP(&jwtSecret, "jwt-secret", "", "", "JWT secret key (overrides config)")
	serverCmd.Flags().IntVarP(&shutdownTimeout, "shutdown-timeout", "", 0, "Server shutdown timeout in seconds (overrides config)")

	// Bind flags to viper
	viper.BindPFlag("config", serverCmd.Flags().Lookup("config"))
	viper.BindPFlag("env", serverCmd.Flags().Lookup("env"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.host", serverCmd.Flags().Lookup("host"))
	viper.BindPFlag("debug", serverCmd.Flags().Lookup("debug"))
	viper.BindPFlag("logger.level", serverCmd.Flags().Lookup("log-level"))
	viper.BindPFlag("logger.format", serverCmd.Flags().Lookup("log-format"))
	viper.BindPFlag("database.uri", serverCmd.Flags().Lookup("db-uri"))
	viper.BindPFlag("database.name", serverCmd.Flags().Lookup("db-name"))
	viper.BindPFlag("auth.secret", serverCmd.Flags().Lookup("jwt-secret"))
	viper.BindPFlag("server.shutdown_timeout", serverCmd.Flags().Lookup("shutdown-timeout"))
}

func startServer() {
	// Load configuration with flags
	cfg, err := loadConfigWithFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewZapLogger(cfg.Logger.Development)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Log configuration info
	log.InfoLegacy("Starting server with configuration",
		logger.String("config_path", configPath),
		logger.String("environment", environment),
		logger.String("host", cfg.Server.Host),
		logger.Int("port", cfg.Server.Port),
		logger.Bool("debug", debug),
		logger.String("log_level", cfg.Logger.Level),
		logger.String("db_name", cfg.Database.Database))

	// Initialize MongoDB
	mongoDB, err := database.NewMongoDB(&cfg.Database, log)
	if err != nil {
		log.FatalLegacy("Failed to connect to MongoDB", logger.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoDB.Close(ctx); err != nil {
			log.ErrorLegacy("Failed to close MongoDB connection", logger.Error(err))
		}
	}()

	// Initialize repositories
	userRepo := mongodb.NewUserRepository(mongoDB.Database)
	foodRepo := mongodb.NewFoodRepository(mongoDB.Database)
	mealTemplateRepo := mongodb.NewMealTemplateRepository(mongoDB.Database)
	mealPlanRepo := mongodb.NewMealPlanRepository(mongoDB.Database)
	shoppingRepo := mongodb.NewShoppingListRepository(mongoDB.Database)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth, log)
	foodService := service.NewFoodService(foodRepo, log)
	mealService := service.NewMealService(mealTemplateRepo, foodRepo, log)
	mealPlanService := service.NewMealPlanService(mealPlanRepo, mealTemplateRepo, log)
	shoppingService := service.NewShoppingService(shoppingRepo, mealPlanRepo, log)
	reportService := service.NewReportService(mealPlanRepo, log)

	// Initialize handlers
	handlers := rest.NewHandlers(
		authService,
		foodService,
		mealService,
		mealPlanService,
		shoppingService,
		reportService,
		mongoDB.Client,
		log,
	)

	// Setup Gin router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Setup routes
	rest.SetupRoutes(router, handlers)

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.InfoLegacy("Server started successfully",
			logger.String("address", server.Addr),
			logger.String("mode", cfg.Server.Mode))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.FatalLegacy("Failed to start server", logger.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.InfoLegacy("Shutting down server...")

	// Shutdown server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.ErrorLegacy("Server forced to shutdown", logger.Error(err))
	}

	log.InfoLegacy("Server exited")
}

// loadConfigWithFlags loads configuration with command line flags override
func loadConfigWithFlags() (*config.Config, error) {
	// Set default values
	viper.SetDefault("env", "dev")
	viper.SetDefault("debug", false)
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "console")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.shutdown_timeout", 30)

	// Determine config file path
	if configPath == "" {
		env := environment
		if env == "" {
			env = viper.GetString("env")
		}
		configPath = fmt.Sprintf("./configs/config.%s.yaml", env)
	}

	// Load config file
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	viper.AutomaticEnv()

	// Override with flags
	if port > 0 {
		viper.Set("server.port", port)
	}
	if host != "" {
		viper.Set("server.host", host)
	}
	if debug {
		viper.Set("debug", true)
		viper.Set("logger.level", "debug")
		viper.Set("server.mode", "debug")
	}
	if logLevel != "" {
		viper.Set("logger.level", logLevel)
	}
	if logFormat != "" {
		viper.Set("logger.format", logFormat)
	}
	if dbURI != "" {
		viper.Set("database.uri", dbURI)
	}
	if dbName != "" {
		viper.Set("database.name", dbName)
	}
	if jwtSecret != "" {
		viper.Set("auth.secret", jwtSecret)
	}
	if shutdownTimeout > 0 {
		viper.Set("server.shutdown_timeout", shutdownTimeout)
	}

	// Unmarshal config
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set logger development mode based on debug flag or environment
	cfg.Logger.Development = debug || cfg.Logger.Level == "debug"

	return &cfg, nil
}
