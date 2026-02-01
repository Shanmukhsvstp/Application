package routes

import (
	"application/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(app *fiber.App) {

	api := app.Group("/api")
	auth := api.Group("/auth")

	authHandler := &handlers.AuthHandler{}

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	auth.Post("/login", authHandler.Login)
}
