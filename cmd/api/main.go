package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/entorno35/backend/internal/adapters/http"
	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/config"
	"github.com/entorno35/backend/internal/core/jwt"
	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/database"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	dsn, err := cfg.DatabaseURL()
	if err != nil {
		log.Fatalf("Failed to get database URL: %v", err)
	}

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize services
	jwtService := jwt.NewService(cfg.JWT.Secret)

	// Parse token expiry from config (default: 24h)
	tokenExpiry := 24 * time.Hour
	if cfg.JWT.Expiry != "" {
		parsedExpiry, err := time.ParseDuration(cfg.JWT.Expiry)
		if err == nil {
			tokenExpiry = parsedExpiry
		}
	}

	// Initialize repositories
	authRepo := postgres.NewAuthRepository(db)
	assessmentRepo := postgres.NewAssessmentRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	staffRepo := postgres.NewStaffRepository(db)

	// Initialize services
	scoringService := services.NewScoringService(assessmentRepo)
	assessmentService := services.NewAssessmentService(assessmentRepo, companyRepo, staffRepo)

	// Initialize handlers
	authHandler := http.NewAuthHandler(jwtService, authRepo, tokenExpiry)
	scoringHandler := http.NewScoringHandler(scoringService)
	assessmentHandler := http.NewAssessmentHandler(assessmentService)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "entorno35-api",
		})
	})

	// Auth endpoints
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Protected API endpoints (require authentication)
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtService))
	api.Use(middleware.TenantMiddleware())
	{
		// Assessment endpoints
		assessments := api.Group("/assessments")
		{
			assessments.POST("", assessmentHandler.CreateAssessment)
			assessments.GET("", assessmentHandler.ListAssessments)
			assessments.GET("/:id", assessmentHandler.GetAssessment)
			assessments.POST("/:id/links", assessmentHandler.CreateAssessmentLink)
			assessments.POST("/:id/calculate", scoringHandler.CalculateAssessment)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

