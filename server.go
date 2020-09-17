package main

import (
	"log"

	"github.com/duyet/applause-btn/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := Setup()

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Setup Setup a fiber app with all of its routes
func Setup() *fiber.App {
	app := fiber.New()

	// CORS
	app.Use(cors.New())
	// Compression config
	app.Use(compress.New())
	// Logger
	app.Use(logger.New())

	// Setup routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Hello, World!"))
	})
	app.Get("/get-claps", api.GetClaps)
	app.Post("/get-multiple", api.GetMultiple)
	app.Post("/update-claps", api.UpdateClaps)
	app.Static("/public", "./public")

	return app
}
