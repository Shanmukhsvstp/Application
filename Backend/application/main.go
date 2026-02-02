package main

import (
	"application/db"
	"application/migrate"
	"application/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	// # DATABASE_URL="postgresql://postgres:mmffTqNJJ9l2qgdb@db.tttefpaywxxrlwfqxcgu.supabase.co:5432/postgres"
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUrl := os.Getenv("DATABASE_URL")

	dbPool := db.Connect(dbUrl)
	defer dbPool.Close()

	routes.SetupApiRoutes(app, dbPool)

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		migrate.RunMigrations(dbPool)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Fiber")
	})

	app.Listen(":3000")
}
