package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"nutrient_be/internal/config"
	"nutrient_be/internal/database"
	"nutrient_be/internal/pkg/logger"
)

// Migrate flags
var (
	migrateConfigPath  string
	migrateEnvironment string
	migrateDbURI       string
	migrateDbName      string
	migrateDryRun      bool
	migrateForce       bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long: `Run database migrations to create indexes and setup collections.

Examples:
  # Run migrations with default config
  nutrient-api migrate

  # Run migrations with custom config
  nutrient-api migrate --config=./configs/config.prod.yaml

  # Run migrations with custom database
  nutrient-api migrate --db-uri=mongodb://localhost:27017 --db-name=nutrient_prod

  # Dry run (show what would be done without executing)
  nutrient-api migrate --dry-run

  # Force run (skip confirmation prompts)
  nutrient-api migrate --force`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Migration flags
	migrateCmd.Flags().StringVarP(&migrateConfigPath, "config", "c", "", "Path to configuration file")
	migrateCmd.Flags().StringVarP(&migrateEnvironment, "env", "e", "", "Environment (dev, staging, prod)")
	migrateCmd.Flags().StringVarP(&migrateDbURI, "db-uri", "", "", "MongoDB connection URI (overrides config)")
	migrateCmd.Flags().StringVarP(&migrateDbName, "db-name", "", "", "MongoDB database name (overrides config)")
	migrateCmd.Flags().BoolVarP(&migrateDryRun, "dry-run", "", false, "Show what would be done without executing")
	migrateCmd.Flags().BoolVarP(&migrateForce, "force", "f", false, "Skip confirmation prompts")

	// Bind flags to viper
	viper.BindPFlag("migrate.config", migrateCmd.Flags().Lookup("config"))
	viper.BindPFlag("migrate.env", migrateCmd.Flags().Lookup("env"))
	viper.BindPFlag("migrate.db_uri", migrateCmd.Flags().Lookup("db-uri"))
	viper.BindPFlag("migrate.db_name", migrateCmd.Flags().Lookup("db-name"))
	viper.BindPFlag("migrate.dry_run", migrateCmd.Flags().Lookup("dry-run"))
	viper.BindPFlag("migrate.force", migrateCmd.Flags().Lookup("force"))
	viper.BindPFlag("migrate.force", migrateCmd.Flags().Lookup("force"))
}

func runMigrations() {
	// Load configuration with flags
	cfg, err := loadMigrateConfig()
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

	// Log migration info
	log.Info(context.Background(), "Starting database migrations",
		logger.String("config_path", migrateConfigPath),
		logger.String("environment", migrateEnvironment),
		logger.String("db_name", cfg.Database.Database),
		logger.Bool("dry_run", migrateDryRun),
		logger.Bool("force", migrateForce))

	// Confirm migration if not forced
	if !migrateForce && !migrateDryRun {
		fmt.Printf("Are you sure you want to run migrations on database '%s'? (y/N): ", cfg.Database.Database)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Migration cancelled")
			return
		}
	}

	// Initialize MongoDB
	mongoDB, err := database.NewMongoDB(&cfg.Database, log)
	if err != nil {
		log.Fatal(context.Background(), "Failed to connect to MongoDB", logger.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoDB.Close(ctx); err != nil {
			log.Error(context.Background(), "Failed to close MongoDB connection", logger.Error(err))
		}
	}()

	// Run migrations
	if err := runDatabaseMigrations(mongoDB, log); err != nil {
		log.Fatal(context.Background(), "Failed to run migrations", logger.Error(err))
	}

	log.Info(context.Background(), "Migrations completed successfully")
}

func loadMigrateConfig() (*config.Config, error) {
	// Set default values
	viper.SetDefault("env", "dev")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "console")

	// Determine config file path
	if migrateConfigPath == "" {
		env := migrateEnvironment
		if env == "" {
			env = viper.GetString("env")
		}
		migrateConfigPath = fmt.Sprintf("./configs/config.%s.yaml", env)
	}

	// Load config file
	viper.SetConfigFile(migrateConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	viper.AutomaticEnv()

	// Override with flags
	if migrateDbURI != "" {
		viper.Set("database.uri", migrateDbURI)
	}
	if migrateDbName != "" {
		viper.Set("database.name", migrateDbName)
	}

	// Unmarshal config
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func runDatabaseMigrations(mongoDB *database.MongoDB, log logger.Logger) error {
	// Create indexes
	collections := []string{"users", "foods", "meal_templates", "meal_plans", "shopping_lists"}

	for _, collectionName := range collections {
		collection := mongoDB.GetCollection(collectionName)

		if migrateDryRun {
			log.Info(context.Background(), "DRY RUN: Would create indexes for collection", logger.String("collection", collectionName))
			continue
		}

		// Create indexes based on collection
		switch collectionName {
		case "users":
			_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    bson.M{"email": 1},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				return fmt.Errorf("failed to create users index: %w", err)
			}
		case "foods":
			indexes := []mongo.IndexModel{
				{Keys: bson.M{"searchTerms": "text"}},
				{Keys: bson.M{"createdBy": 1, "visibility": 1}},
				{Keys: bson.M{"category": 1}},
				{Keys: bson.M{"source": 1}},
			}
			_, err := collection.Indexes().CreateMany(context.Background(), indexes)
			if err != nil {
				return fmt.Errorf("failed to create foods indexes: %w", err)
			}
		case "meal_templates":
			indexes := []mongo.IndexModel{
				{Keys: bson.M{"userId": 1, "mealType": 1}},
				{Keys: bson.M{"userId": 1, "isPublic": 1}},
			}
			_, err := collection.Indexes().CreateMany(context.Background(), indexes)
			if err != nil {
				return fmt.Errorf("failed to create meal_templates indexes: %w", err)
			}
		case "meal_plans":
			indexes := []mongo.IndexModel{
				{Keys: bson.M{"userId": 1, "startDate": -1}},
				{Keys: bson.M{"userId": 1, "planType": 1}},
				{Keys: bson.M{"userId": 1, "status": 1}},
			}
			_, err := collection.Indexes().CreateMany(context.Background(), indexes)
			if err != nil {
				return fmt.Errorf("failed to create meal_plans indexes: %w", err)
			}
		case "shopping_lists":
			indexes := []mongo.IndexModel{
				{Keys: bson.M{"userId": 1, "mealPlanId": 1}},
				{Keys: bson.M{"userId": 1, "status": 1}},
			}
			_, err := collection.Indexes().CreateMany(context.Background(), indexes)
			if err != nil {
				return fmt.Errorf("failed to create shopping_lists indexes: %w", err)
			}
		}

		log.Info(context.Background(), "Created indexes for collection", logger.String("collection", collectionName))
	}

	return nil
}
