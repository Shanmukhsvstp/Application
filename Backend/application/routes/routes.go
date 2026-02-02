package routes

import (
	"application/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupApiRoutes(app *fiber.App, dbPool *pgxpool.Pool) {

	api := app.Group("/api")
	auth := api.Group("/auth")

	authHandler := handlers.DBHandler(dbPool)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	auth.Post("/login", authHandler.Login)
	auth.Post("/signup", authHandler.Signup)
	auth.Get("/validate", authHandler.ValidateUserToken)
}
