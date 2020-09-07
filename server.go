package main

import (
	"log"

	"github.com/duyet/applause-btn/api"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func main() {
	app := Setup()
	app = SetupMiddleware(app)

	// Start server
	log.Fatal(app.Listen(3000))
}

// Setup Setup a fiber app with all of its routes
func Setup() *fiber.App {
	app := fiber.New()

	// Setup routes
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, World!")
	})

	app.Get("/get-claps", api.GetClaps)
	app.Post("/get-multiple", api.GetMultiple)
	app.Post("/update-claps", api.UpdateClaps)
	app.Static("/public", "./public")

	return app
}

// SetupMiddleware setup middleware for app
func SetupMiddleware(app *fiber.App) *fiber.App {
	// CORS
	app.Use(cors.New())
	// Compression config
	app.Use(middleware.Compress())
	// Logger
	app.Use(middleware.Logger())

	return app
}
