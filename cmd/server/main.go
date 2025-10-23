package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leftathome/trunchbull/internal/api"
	"github.com/leftathome/trunchbull/internal/config"
	"github.com/leftathome/trunchbull/internal/db"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	setupLogging(cfg)

	log.Info("Starting Trunchbull Student Dashboard...")

	// Initialize database
	database, err := db.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Setup API router
	router := setupRouter(cfg, database)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Infof("Server listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogging(cfg *config.Config) {
	// Set log level
	level, err := log.ParseLevel(cfg.Logging.Level)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to info", cfg.Logging.Level)
		level = log.InfoLevel
	}
	log.SetLevel(level)

	// Set log format
	if cfg.Logging.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	log.SetOutput(os.Stdout)
}

func setupRouter(cfg *config.Config, database *db.DB) *gin.Engine {
	// Set gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"version": "0.1.0",
		})
	})

	// API routes
	apiHandler := api.NewHandler(cfg, database)
	apiGroup := router.Group("/api")
	{
		apiHandler.RegisterRoutes(apiGroup)
	}

	// TODO: Serve frontend static files
	// router.Static("/", "./frontend/build")

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
