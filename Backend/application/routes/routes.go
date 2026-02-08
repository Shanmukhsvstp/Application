package routes

import (
	"application/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupApiRoutes(app *fiber.App, dbPool *pgxpool.Pool) {

	api := app.Group("/api")
	auth := api.Group("/auth")
	user := app.Group("/user")

	authHandler := handlers.DBHandler(dbPool)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	auth.Post("/login", authHandler.Login)
	auth.Post("/signup", authHandler.Signup)
	auth.Get("/validate", authHandler.ValidateUserToken)

	// Email Verification Related APIs
	user.Get("/send_verification_email", authHandler.SendVerificationEmail)
	user.Post("/verify_email", authHandler.VerifyEmail)
	user.Get("/is_email_verified", authHandler.IsEmailVerified)
}
