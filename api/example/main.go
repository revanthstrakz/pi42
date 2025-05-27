package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/revanthstrakz/pi42/api"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default error handler
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-API-Key, X-API-Secret",
	}))

	// Routes
	api.SetupPi42Routes(app, createAuthMiddleware())

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// Simple auth middleware example - customize as needed
func createAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for API key and secret in headers or auth token
		apiKey := c.Get("X-API-Key")
		apiSecret := c.Get("X-API-Secret")

		if apiKey == "" || apiSecret == "" {
			// Extract from Authorization header if needed
			// Example: Parse JWT token and extract API keys from claims

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API credentials required",
			})
		}

		// Continue to next handler
		return c.Next()
	}
}
