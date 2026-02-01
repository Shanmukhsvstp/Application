package handlers

import (
	"application/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if req.Password != "123456" {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	return c.JSON(models.LoginResponse{
		Token: "abc123",
	})
}
