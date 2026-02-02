package handlers

import (
	"application/models"
	"application/tools"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func DBHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	currentUsername := req.Username
	userPassword := ""

	userId := ""

	err := h.DB.QueryRow(
		c.Context(),
		`
		SELECT id, password FROM users WHERE username=$1 OR email=$1
	`,
		currentUsername,
	).Scan(&userId, &userPassword)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{
				"error": "User not found, please sign up first!",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if !tools.PasswordMatches(req.Password, userPassword) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Incorrect password!",
		})
	}

	// Password check done, email/username is valid, now finalize authentication

	token, err := tools.GenerateToken(userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(models.LoginResponse{
		Token: token,
	})
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *fiber.Ctx) error {

	var req models.SignupRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	currentEmail := req.Email
	currentUsername := req.Username

	if currentUsername == "" || currentEmail == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "username and email are required",
		})
	}

	if len(currentUsername) < 3 || len(currentUsername) > 30 {
		return c.Status(400).JSON(fiber.Map{
			"error": "username must be between 3 and 30 characters",
		})
	}

	if len(req.Password) < 8 || len(req.Password) > 100 {
		return c.Status(400).JSON(fiber.Map{
			"error": "password must be between 8 and 100 characters",
		})
	}

	// Above are checks for lengths, now if user already exists
	exists, err := tools.UserAlreadyExist(currentEmail, h.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	if exists {
		return c.Status(409).JSON(fiber.Map{
			"error": "user with this email already exists",
		})
	}

	// Above are checks for lengths, now username availaibility
	isUnique, err := tools.UsernameIsUnique(currentUsername, h.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	if !isUnique {
		return c.Status(409).JSON(fiber.Map{
			"error": "username already taken",
		})
	}

	// All validations passed, create the user

	// Hash the password
	hashedPassword, err := tools.HashPassword(req.Password)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var userID string

	err = h.DB.QueryRow(
		c.Context(),
		`
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`,
		currentUsername,
		currentEmail,
		hashedPassword,
	).Scan(&userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	token, err := tools.GenerateToken(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(models.SignupResponse{
		Token: token,
	})
}

func (h *AuthHandler) CheckUsername(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "username is available",
	})
}

func (h *AuthHandler) ValidateUserToken(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}
	token := strings.TrimSpace(
		strings.TrimPrefix(authHeader, "Bearer "),
	)

	if token == authHeader {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid authorization format",
		})
	}

	userID, exp, _, err := tools.GetDataFromToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	expTime := int64(exp)
	now := time.Now().Unix()

	if expTime <= now+(24*60*60) { // if token expires in next 24 hours
		newToken, err := tools.GenerateToken(userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to generate new token",
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"token": newToken,
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "token is valid",
	})
}
