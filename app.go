package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// CatFactResponse represents the response from the Cat Facts API
type CatFactResponse struct {
	Fact string `json:"fact"`
}

// ProfileResponse represents the response for the /me endpoint
type ProfileResponse struct {
	Status    string   `json:"status"`
	User      UserInfo `json:"user"`
	Timestamp string   `json:"timestamp"`
	Fact      string   `json:"fact"`
}

// UserInfo contains user details
type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Stack string `json:"stack"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Config holds application configuration
type Config struct {
	Email             string
	Name              string
	Stack             string
	CatFactAPIURL     string
	CatFactAPITimeout time.Duration
	Port              string
}

// loadConfig loads configuration from environment variables or defaults
func loadConfig() *Config {
	return &Config{
		Email:             getEnv("USER_EMAIL", "macmartins081@gmail.com"),
		Name:              getEnv("USER_NAME", "Obinna Aguwa"),
		Stack:             getEnv("USER_STACK", "Go/Gin"),
		CatFactAPIURL:     getEnv("CAT_FACT_API", "https://catfact.ninja/fact"),
		CatFactAPITimeout: 5 * time.Second,
		Port:              getEnv("PORT", "8080"),
	}
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// fetchCatFact fetches a cat fact from the external API with timeout
func fetchCatFact(config *Config) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.CatFactAPITimeout,
	}

	// Fetch the cat fact
	resp, err := client.Get(config.CatFactAPIURL)
	if err != nil {
		log.Printf("Error fetching cat fact: %v", err)
		return "A cat's hearing is the most acute of all domestic animals, being able to hear high-pitched sounds up to 64,000 Hz.", nil
	}
	defer resp.Body.Close()

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Cat Facts API returned status: %d", resp.StatusCode)
		return "A cat's purr vibrates at a frequency that may promote bone and muscle growth.", nil
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading cat fact response: %v", err)
		return "A cat can rotate its ears independently up to 180 degrees.", nil
	}

	// Parse JSON response
	var catResp CatFactResponse
	if err := json.Unmarshal(body, &catResp); err != nil {
		log.Printf("Error parsing cat fact JSON: %v", err)
		return "Cats have a specialized collarbone that allows them to squeeze through tight spaces.", nil
	}

	if catResp.Fact == "" {
		log.Println("Empty cat fact received")
		return "A cat's nose print is unique, much like a human's fingerprint.", nil
	}

	return catResp.Fact, nil
}

// setupRoutes configures all API routes
func setupRoutes(router *gin.Engine, config *Config) {
	// Middleware for CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Profile endpoint
	router.GET("/me", func(c *gin.Context) {
		// Fetch cat fact dynamically on each request
		catFact, err := fetchCatFact(config)
		if err != nil {
			log.Printf("Failed to fetch cat fact: %v", err)
		}

		// Create response with current UTC timestamp
		response := ProfileResponse{
			Status: "success",
			User: UserInfo{
				Email: config.Email,
				Name:  config.Name,
				Stack: config.Stack,
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Fact:      catFact,
		}

		// Log request
		log.Printf("GET /me - Status: 200, User: %s", config.Name)

		// Send response
		c.JSON(http.StatusOK, response)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Welcome endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Obinna Aguwa Profile API",
			"usage":   "GET /me for profile information",
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Status:  "error",
			Message: "Endpoint not found",
		})
	})
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	config := loadConfig()

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Obinna Aguwa Profile API server on port %s", config.Port)
	log.Printf("User: %s (%s) - Stack: %s", config.Name, config.Email, config.Stack)

	// Set Gin mode (release for production, debug for development)
	// Uncomment for production: gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// Create Gin router
	router := gin.Default()

	// Setup routes
	setupRoutes(router, config)

	// Start server
	log.Printf("API server listening at http://localhost:%s", config.Port)
	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
