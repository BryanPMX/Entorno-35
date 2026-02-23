package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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
		&domain.PendingCompanyRegistration{},
		&domain.StripeWebhookEvent{},
		&domain.Staff{},
		&domain.Category{},
		&domain.Domain{},
		&domain.Dimension{},
		&domain.Question{},
		&domain.Assessment{},
		&domain.Response{},
		&domain.AssessmentLink{},
	); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize services
	jwtService := jwt.NewService(cfg.JWT.Secret)
	defaultOrigin := strings.TrimRight(cfg.CORS.Origin, "/")
	stripeSuccessURL := cfg.Stripe.SuccessURL
	if stripeSuccessURL == "" {
		stripeSuccessURL = defaultOrigin + "/register/success?session_id={CHECKOUT_SESSION_ID}"
	}
	stripeCancelURL := cfg.Stripe.CancelURL
	if stripeCancelURL == "" {
		stripeCancelURL = defaultOrigin + "/register?canceled=1"
	}
	stripeService := services.NewStripeService(services.StripeServiceConfig{
		SecretKey:      cfg.Stripe.SecretKey,
		WebhookSecret:  cfg.Stripe.WebhookSecret,
		MonthlyPriceID: cfg.Stripe.MonthlyPriceID,
		YearlyPriceID:  cfg.Stripe.YearlyPriceID,
		SuccessURL:     stripeSuccessURL,
		CancelURL:      stripeCancelURL,
	})

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
	pendingRegistrationRepo := postgres.NewPendingRegistrationRepository(db)
	stripeWebhookEventRepo := postgres.NewStripeWebhookEventRepository(db)
	assessmentRepo := postgres.NewAssessmentRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	staffRepo := postgres.NewStaffRepository(db)
	responseRepo := postgres.NewResponseRepository(db)
	reportRepo := postgres.NewReportRepository(db)

	// Initialize services
	scoringService := services.NewScoringService(assessmentRepo)
	assessmentService := services.NewAssessmentService(assessmentRepo, companyRepo, staffRepo, responseRepo, scoringService)
	staffService := services.NewStaffService(staffRepo)
	reportService := services.NewReportService(reportRepo)

	if cfg.Cleanup.PendingRegistrationEnabled {
		cleanupInterval, err := time.ParseDuration(cfg.Cleanup.PendingRegistrationInterval)
		if err != nil {
			log.Fatalf("Invalid PENDING_REGISTRATION_CLEANUP_INTERVAL '%s': %v", cfg.Cleanup.PendingRegistrationInterval, err)
		}
		cleanupRetention, err := time.ParseDuration(cfg.Cleanup.PendingRegistrationRetention)
		if err != nil {
			log.Fatalf("Invalid PENDING_REGISTRATION_CLEANUP_RETENTION '%s': %v", cfg.Cleanup.PendingRegistrationRetention, err)
		}
		if cleanupInterval <= 0 {
			log.Fatalf("PENDING_REGISTRATION_CLEANUP_INTERVAL must be > 0")
		}
		if cleanupRetention < 0 {
			log.Fatalf("PENDING_REGISTRATION_CLEANUP_RETENTION must be >= 0")
		}

		cleanupService := services.NewPendingRegistrationCleanupService(pendingRegistrationRepo, cleanupInterval, cleanupRetention)
		stopCleanup := cleanupService.Start()
		defer stopCleanup()
		log.Printf("Pending registration cleanup enabled (interval=%s, retention=%s)", cleanupInterval, cleanupRetention)
	} else {
		log.Println("Pending registration cleanup disabled")
	}

	// Initialize handlers
	authHandler := http.NewAuthHandler(jwtService, authRepo, tokenExpiry)
	scoringHandler := http.NewScoringHandler(scoringService)
	assessmentHandler := http.NewAssessmentHandler(assessmentService)
	staffHandler := http.NewStaffHandler(staffService)
	reportHandler := http.NewReportHandler(reportService)
	billingHandler := http.NewBillingHandler(authRepo, pendingRegistrationRepo, stripeWebhookEventRepo, stripeService, stripeSuccessURL, stripeCancelURL)

	// Initialize router with custom middleware (avoid double CORS)
	router := gin.New()

	// Add default middleware manually
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint (registered before CORS to ensure it's always accessible)
	// Supports both GET and HEAD for Docker/Kubernetes health checks
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "entorno35-api",
		})
	})
	router.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// Configure CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORS.Origin}
	corsConfig.AllowMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Auth endpoints
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Billing endpoints (public for checkout + webhook)
	billing := router.Group("/billing")
	{
		billing.POST("/checkout-session", billingHandler.CreateCheckoutSession)
		billing.POST("/webhook", billingHandler.HandleWebhook)
		billing.GET("/checkout-session/:id/verify", billingHandler.VerifyCheckoutSession)
	}

	// Public API endpoints (no authentication required)
	public := router.Group("/api/v1/assessments/public")
	{
		public.GET("/:token", assessmentHandler.GetPublicAssessment)
		public.POST("/:token/submit", assessmentHandler.SubmitAssessment)
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
