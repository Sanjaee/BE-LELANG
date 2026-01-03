package app

import (
	"log"
	"time"
	"yourapp/internal/config"
	"yourapp/internal/middleware"
	"yourapp/internal/model"
	"yourapp/internal/repository"
	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.ServerPort == "5000" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(corsMiddleware(cfg.ClientURL))

	// Rate limiting middleware (if enabled)
	if cfg.RateLimitEnabled {
		rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
		r.Use(rateLimiter.Middleware())
		log.Printf("Rate limiting enabled: %d req/sec, burst: %d", cfg.RateLimitRPS, cfg.RateLimitBurst)
	}

	// Initialize database
	db, err := initDB(cfg)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Auto migrate all models
	if err := db.AutoMigrate(
		&model.User{},
		&model.Seller{},
		&model.Organizer{},
		&model.ItemCategory{},
		&model.AuctionItem{},
		&model.ItemImage{},
		&model.AuctionSchedule{},
		&model.Bid{},
	); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sellerRepo := repository.NewSellerRepository(db)
	organizerRepo := repository.NewOrganizerRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	itemRepo := repository.NewAuctionItemRepository(db)
	imageRepo := repository.NewItemImageRepository(db)
	scheduleRepo := repository.NewAuctionScheduleRepository(db)
	bidRepo := repository.NewBidRepository(db)

	// Initialize RabbitMQ with retry logic
	rabbitMQ := initRabbitMQWithRetry(cfg)

	// Initialize email service
	emailService := service.NewEmailService(cfg)

	// Initialize email worker if RabbitMQ is available
	var emailWorker *service.EmailWorker
	if rabbitMQ != nil {
		emailWorker = service.NewEmailWorker(emailService, rabbitMQ)
		if err := emailWorker.Start(); err != nil {
			log.Printf("Warning: Failed to start email worker: %v", err)
		} else {
			log.Println("Email worker started successfully")
		}
	} else {
		log.Println("Email worker not started - RabbitMQ connection failed. Will retry on first email send.")
		// Start background goroutine to retry RabbitMQ connection and start email worker
		go func() {
			for {
				time.Sleep(10 * time.Second)
				newRabbitMQ := initRabbitMQWithRetry(cfg)
				if newRabbitMQ != nil {
					log.Println("RabbitMQ reconnected! Starting email worker...")
					emailWorker = service.NewEmailWorker(emailService, newRabbitMQ)
					if err := emailWorker.Start(); err != nil {
						log.Printf("Warning: Failed to start email worker after reconnect: %v", err)
					} else {
						log.Println("Email worker started successfully after reconnect")
						break
					}
				}
			}
		}()
	}

	// Initialize services
	authService := service.NewAuthServiceWithConfig(userRepo, cfg.JWTSecret, rabbitMQ, cfg)
	auctionService := service.NewAuctionService(
		sellerRepo,
		organizerRepo,
		categoryRepo,
		itemRepo,
		imageRepo,
		scheduleRepo,
		bidRepo,
		userRepo,
	)

	// Initialize handlers
	authHandler := NewAuthHandler(authService, cfg.JWTSecret)
	auctionHandler := NewAuctionHandler(auctionService, cfg.JWTSecret)

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/resend-otp", authHandler.ResendOTP)
			auth.POST("/google-oauth", authHandler.GoogleOAuth)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.RequestResetPassword)
			auth.POST("/verify-reset-password", authHandler.VerifyResetPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-email", authHandler.VerifyEmail)

			// Protected routes
			auth.GET("/me", authHandler.AuthMiddleware(), authHandler.GetMe)
		}

		// Auction routes (public)
		auctions := api.Group("/auctions")
		{
			// Public endpoints for frontend
			auctions.GET("", auctionHandler.GetAuctionItemsForFrontend)
			auctions.GET("/:id", auctionHandler.GetAuctionItem)
			auctions.GET("/:id/bids", auctionHandler.GetItemBids)

			// Categories
			auctions.GET("/categories", auctionHandler.GetCategories)
			auctions.GET("/categories/:id", auctionHandler.GetCategory)
		}

		// Admin auction management (protected)
		adminAuctions := api.Group("/admin/auctions")
		adminAuctions.Use(authHandler.AuthMiddleware())
		{
			// Sellers
			adminAuctions.POST("/sellers", auctionHandler.CreateSeller)
			adminAuctions.GET("/sellers", auctionHandler.GetSellers)
			adminAuctions.GET("/sellers/:id", auctionHandler.GetSeller)

			// Organizers
			adminAuctions.POST("/organizers", auctionHandler.CreateOrganizer)
			adminAuctions.GET("/organizers", auctionHandler.GetOrganizers)
			adminAuctions.GET("/organizers/:id", auctionHandler.GetOrganizer)

			// Categories
			adminAuctions.POST("/categories", auctionHandler.CreateCategory)

			// Items
			adminAuctions.POST("/items", auctionHandler.CreateAuctionItem)
			adminAuctions.GET("/items", auctionHandler.GetAuctionItems)
			adminAuctions.PUT("/items/:id", auctionHandler.UpdateAuctionItem)
			adminAuctions.POST("/items/:id/publish", auctionHandler.PublishAuctionItem)
			adminAuctions.DELETE("/items/:id", auctionHandler.DeleteAuctionItem)
		}

		// Bidding routes (protected)
		bids := api.Group("/bids")
		bids.Use(authHandler.AuthMiddleware())
		{
			bids.POST("", auctionHandler.PlaceBid)
			bids.GET("/my-bids", auctionHandler.GetUserBids)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	if dsn == "" {
		dsn = "host=" + cfg.PostgresHost +
			" port=" + cfg.PostgresPort +
			" user=" + cfg.PostgresUser +
			" password=" + cfg.PostgresPassword +
			" dbname=" + cfg.PostgresDB +
			" sslmode=" + cfg.PostgresSSLMode
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// initRabbitMQWithRetry attempts to connect to RabbitMQ with exponential backoff retry
func initRabbitMQWithRetry(cfg *config.Config) *util.RabbitMQClient {
	maxRetries := 10
	initialDelay := 2 * time.Second
	maxDelay := 30 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		rabbitMQ, err := util.NewRabbitMQClient(cfg)
		if err == nil {
			log.Printf("RabbitMQ connected successfully on attempt %d", attempt)
			return rabbitMQ
		}

		if attempt < maxRetries {
			// Calculate delay with exponential backoff
			delay := initialDelay * time.Duration(1<<uint(attempt-1))
			if delay > maxDelay {
				delay = maxDelay
			}

			log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, delay)
			time.Sleep(delay)
		} else {
			log.Printf("Warning: Failed to connect to RabbitMQ after %d attempts: %v. Email sending will be disabled.", maxRetries, err)
			log.Println("Note: RabbitMQ will be retried automatically when email is sent (if connection is restored)")
		}
	}

	return nil
}

func corsMiddleware(clientURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", clientURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
