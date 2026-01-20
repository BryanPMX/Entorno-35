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
	"github.com/entorno35/backend/internal/domain"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-contrib/cors"
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
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			log.Printf("Warning: failed to close database connection: %v", closeErr)
		}
		// Close audit service to ensure all logs are flushed
		if closeErr := services.GetAuditService().Close(); closeErr != nil {
			log.Printf("Warning: failed to close audit service: %v", closeErr)
		}
	}()

	// Run database migrations (create/update schema)
	log.Println("Running database migrations...")
	if err := database.Migrate(
		&domain.Company{},
		&domain.Staff{},
		&domain.Category{},
		&domain.Domain{},
		&domain.Dimension{},
		&domain.Question{},
		&domain.Assessment{},
		&domain.Response{},
		&domain.AssessmentLink{},
		&domain.Payment{},
		&domain.Subscription{},
	); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize services
	jwtService := jwt.NewService(cfg.JWT.Secret)

	// Parse token expiry from config (default: 24h)
	tokenExpiry := 24 * time.Hour
	if cfg.JWT.Expiry != "" {
		parsedExpiry, err := time.ParseDuration(cfg.JWT.Expiry)
		if err != nil {
			log.Fatalf("Invalid JWT_EXPIRY duration '%s': %v. Use Go duration format (e.g., '24h', '7200s', '2h30m')", cfg.JWT.Expiry, err)
		}
		tokenExpiry = parsedExpiry
	}

	// Initialize repositories
	authRepo := postgres.NewAuthRepository(db)
	assessmentRepo := postgres.NewAssessmentRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	staffRepo := postgres.NewStaffRepository(db)
	responseRepo := postgres.NewResponseRepository(db)
	reportRepo := postgres.NewReportRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	// Initialize services
	scoringService := services.NewScoringService(assessmentRepo)
	assessmentService := services.NewAssessmentService(assessmentRepo, companyRepo, staffRepo, responseRepo, scoringService)
	staffService := services.NewStaffService(staffRepo)
	reportService := services.NewReportService(reportRepo)
	paymentService := services.NewPaymentService(paymentRepo, cfg.Stripe)

	// Initialize handlers
	authHandler := http.NewAuthHandler(jwtService, authRepo, tokenExpiry)
	scoringHandler := http.NewScoringHandler(scoringService)
	assessmentHandler := http.NewAssessmentHandler(assessmentService)
	staffHandler := http.NewStaffHandler(staffService)
	reportHandler := http.NewReportHandler(reportService)
	paymentHandler := http.NewPaymentHandler(paymentService)

	// Initialize router with custom middleware (avoid double CORS)
	router := gin.New()

	// Add default middleware manually
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORS.Origin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

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

	// Public API endpoints (no authentication required)
	public := router.Group("/api/v1")
	{
		// Assessment endpoints
		publicAssessments := public.Group("/assessments/public")
		{
			publicAssessments.GET("/:token", assessmentHandler.GetPublicAssessment)
			publicAssessments.POST("/:token/submit", assessmentHandler.SubmitAssessment)
		}

		// Payment webhook endpoint
		public.POST("/payments/webhook", paymentHandler.HandleWebhook)
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
			assessments.DELETE("/:id", assessmentHandler.DeleteAssessment)
			assessments.POST("/:id/links", assessmentHandler.CreateAssessmentLink)
			assessments.POST("/:id/send-email", assessmentHandler.SendAssessmentEmail)
			assessments.POST("/:id/calculate", scoringHandler.CalculateAssessment)
		}

		// Staff endpoints
		staff := api.Group("/staff")
		{
			staff.POST("", staffHandler.CreateStaff)
			staff.GET("", staffHandler.ListStaff)
			staff.GET("/:id", staffHandler.GetStaff)
			staff.PUT("/:id", staffHandler.UpdateStaff)
			staff.DELETE("/:id", staffHandler.DeleteStaff)
			staff.POST("/import", staffHandler.ImportStaff)
			staff.POST("/csv/analyze", staffHandler.AnalyzeCSV)
			staff.POST("/csv/preview", staffHandler.PreviewCSVImport)
		}

		// Report endpoints
		reports := api.Group("/reports")
		{
			reports.GET("/individual/:assessment_id", reportHandler.GetIndividualReport)
			reports.GET("/individual/:assessment_id/pdf", reportHandler.GetIndividualReportPDF)
			reports.GET("/general", reportHandler.GetGeneralReport)
			reports.GET("/general/pdf", reportHandler.GetGeneralReportPDF)
		}

		// Payment endpoints
		payments := api.Group("/payments")
		{
			payments.POST("/subscription", paymentHandler.CreateSubscription)
			payments.GET("/subscription", paymentHandler.GetSubscription)
			payments.DELETE("/subscription", paymentHandler.CancelSubscription)
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
