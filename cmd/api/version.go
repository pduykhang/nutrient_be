package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Show version information for the nutrient backend API.

This command displays:
- Application version
- Go version
- Build information
- Git commit hash (if available)`,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show application information",
	Long: `Show detailed application information including:
- Configuration details
- Database connection info
- Environment settings
- Feature flags`,
	Run: func(cmd *cobra.Command, args []string) {
		showInfo()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(infoCmd)
}

func showVersion() {
	fmt.Println("Nutrient Backend API")
	fmt.Println("===================")
	fmt.Printf("Version:     %s\n", getVersion())
	fmt.Printf("Go Version:  %s\n", runtime.Version())
	fmt.Printf("OS/Arch:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Build Time:  %s\n", getBuildTime())
	fmt.Printf("Git Commit:  %s\n", getGitCommit())
	fmt.Println()
}

func showInfo() {
	fmt.Println("Nutrient Backend API - Application Information")
	fmt.Println("=============================================")
	fmt.Println()

	fmt.Println("Configuration:")
	fmt.Printf("  Config Path: %s\n", getConfigPath())
	fmt.Printf("  Environment: %s\n", getEnvironment())
	fmt.Println()

	fmt.Println("Server:")
	fmt.Printf("  Host: %s\n", getServerHost())
	fmt.Printf("  Port: %d\n", getServerPort())
	fmt.Printf("  Mode: %s\n", getServerMode())
	fmt.Println()

	fmt.Println("Database:")
	fmt.Printf("  URI:  %s\n", getDatabaseURI())
	fmt.Printf("  Name: %s\n", getDatabaseName())
	fmt.Println()

	fmt.Println("Features:")
	fmt.Printf("  Context Logging:    Enabled\n")
	fmt.Printf("  Response Middleware: Enabled\n")
	fmt.Printf("  JWT Authentication: Enabled\n")
	fmt.Printf("  MongoDB Support:    Enabled\n")
	fmt.Printf("  Excel Import:       Enabled\n")
	fmt.Println()
}

// Version information - these would typically be set at build time
func getVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}
	return "1.0.0-dev"
}

func getBuildTime() string {
	if buildTime := os.Getenv("BUILD_TIME"); buildTime != "" {
		return buildTime
	}
	return "Unknown"
}

func getGitCommit() string {
	if gitCommit := os.Getenv("GIT_COMMIT"); gitCommit != "" {
		return gitCommit
	}
	return "Unknown"
}

func getConfigPath() string {
	if configPath != "" {
		return configPath
	}
	return "./configs/config.dev.yaml"
}

func getEnvironment() string {
	if environment != "" {
		return environment
	}
	return "dev"
}

func getServerHost() string {
	if host != "" {
		return host
	}
	return "localhost"
}

func getServerPort() int {
	if port > 0 {
		return port
	}
	return 8080
}

func getServerMode() string {
	if debug {
		return "debug"
	}
	return "release"
}

func getDatabaseURI() string {
	if dbURI != "" {
		return dbURI
	}
	return "mongodb://localhost:27017"
}

func getDatabaseName() string {
	if dbName != "" {
		return dbName
	}
	return "nutrient_dev"
}
