package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/duyet/applause-btn/api"
	"github.com/duyet/applause-btn/config"
	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Applause Button service...")
	log.Printf("Configuration: Port=%s, DB=%s", cfg.Port, cfg.DBLocation)

	// Initialize database
	db, err := utils.NewDatabase(cfg.DBLocation)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set global DB for backward compatibility with existing handlers
	// TODO: Refactor all handlers to use dependency injection instead
	if utils.GetDB() == nil {
		utils.DB = db.GetRawDB()
	}

	// Setup Fiber app
	app := Setup(cfg)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("Server listening on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down gracefully...")

	// Shutdown server
	if err := app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close database
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exited")
}

// Setup creates and configures a Fiber app with all routes and middleware
func Setup(cfg *config.Config) *fiber.App {
	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		AppName:               "Applause Button",
		ServerHeader:          "Applause",
		DisableStartupMessage: true,
		ErrorHandler:          customErrorHandler,
	})

	// Recover from panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS configuration
	if cfg != nil && len(cfg.AllowedOrigins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins: joinOrigins(cfg.AllowedOrigins),
			AllowHeaders: "Origin, Content-Type, Accept, Referer",
			AllowMethods: "GET, POST, OPTIONS",
		}))
	} else {
		// Default: allow all origins
		app.Use(cors.New())
	}

	// Compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Request logging
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Health check endpoint
	app.Get("/health", healthCheck)

	// API routes
	app.Get("/", indexHandler)
	app.Get("/get-claps", api.GetClaps)
	app.Get("/get-clappers", api.GetClappers)
	app.Post("/get-multiple", api.GetMultiple)
	app.Post("/update-claps", api.UpdateClaps)

	// Static files
	app.Static("/public", "./public", fiber.Static{
		Compress:      true,
		CacheDuration: 24 * time.Hour,
	})

	return app
}

type indexEndpoints struct {
	Health      string `json:"health"`
	GetClaps    string `json:"get_claps"`
	GetClappers string `json:"get_clappers"`
	GetMultiple string `json:"get_multiple"`
	UpdateClaps string `json:"update_claps"`
}

type indexResponse struct {
	Service   string         `json:"service"`
	Version   string         `json:"version"`
	Status    string         `json:"status"`
	Endpoints indexEndpoints `json:"endpoints"`
}

// indexHandler handles the root endpoint
func indexHandler(c *fiber.Ctx) error {
	return c.JSON(indexResponse{
		Service: "Applause Button",
		Version: "2.0.0",
		Status:  "running",
		Endpoints: indexEndpoints{
			Health:      "/health",
			GetClaps:    "/get-claps",
			GetClappers: "/get-clappers",
			GetMultiple: "/get-multiple (POST)",
			UpdateClaps: "/update-claps (POST)",
		},
	})
}

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

// healthCheck endpoint for monitoring
func healthCheck(c *fiber.Ctx) error {
	// TODO: Add database health check
	return c.JSON(healthResponse{
		Status: "healthy",
		Time:   time.Now().Format(time.RFC3339),
	})
}

type errorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// customErrorHandler provides consistent error responses
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Handle Fiber errors
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log the error
	log.Printf("Error: %v (path: %s, ip: %s)", err, c.Path(), c.IP())

	// Return error response
	return c.Status(code).JSON(errorResponse{
		Error:  err.Error(),
		Status: code,
	})
}

// joinOrigins joins allowed origins with commas
func joinOrigins(origins []string) string {
	if len(origins) == 0 {
		return "*"
	}
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ", "
		}
		result += origin
	}
	return result
}
