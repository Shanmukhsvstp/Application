package handlers

import (
	"application/models"
	"application/tools"
	"fmt"
	"strconv"
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
	var isVerified bool

	err = h.DB.QueryRow(
		c.Context(),
		`SELECT is_verified FROM users WHERE id = $1`,
		userID,
	).Scan(&isVerified)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch user",
		})
	}

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
		"message":     "token is valid",
		"is_verified": isVerified,
	})
}

func (h *AuthHandler) IsEmailVerified(c *fiber.Ctx) error {
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

	userID, _, _, err := tools.GetDataFromToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	var isVerified bool
	err = h.DB.QueryRow(
		c.Context(),
		`SELECT is_verified FROM users WHERE id = $1`,
		userID,
	).Scan(&isVerified)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch user verification status",
		})
	}

	return c.JSON(fiber.Map{
		"is_verified": isVerified,
	})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {

	var req models.VerifyEmailStruct

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	givenOtp := req.Otp

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
	userID, _, _, err := tools.GetDataFromToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	if len(givenOtp) == 6 {
		var storedOtp string
		err = h.DB.QueryRow(
			c.Context(),
			`SELECT otp FROM email_verification_codes WHERE user_id = $1 AND expires_at > NOW()`,
			userID,
		).Scan(&storedOtp)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid or expired OTP",
			})
		}
		if givenOtp != storedOtp {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid OTP",
			})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "OTP must be 6 characters long",
		})
	}

	_, err = h.DB.Exec(
		c.Context(),
		`UPDATE users SET is_verified = true WHERE id = $1`,
		userID,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to verify email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "email verified successfully",
	})

}

func (h *AuthHandler) SendVerificationEmail(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	// isResendReq := c.Query("resend", "false")
	resend, err := strconv.ParseBool(c.Query("resend", "false"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid resend value",
		})
	}

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
	userID, _, _, err := tools.GetDataFromToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	var email string
	var isVerified bool
	err = h.DB.QueryRow(
		c.Context(),
		`SELECT email, is_verified FROM users WHERE id = $1`,
		userID,
	).Scan(&email, &isVerified)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch user email",
		})
	}

	if isVerified {
		return c.Status(400).JSON(fiber.Map{
			"error": "email is already verified",
		})
	}

	// Check if an unexpired otp already exists for the user, if yes then resend the same otp, if not generate a new one and send it
	var existingOTP string
	var expiresAt time.Time
	err = h.DB.QueryRow(
		c.Context(),
		`SELECT otp, expires_at FROM email_verification_codes WHERE user_id = $1 AND expires_at > NOW()`,
		userID,
	).Scan(&existingOTP, &expiresAt)

	if err != nil && err != pgx.ErrNoRows {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to check existing OTP",
		})
	}

	fmt.Printf("DEBUG: existingOTP='%s', resend=%v\n", existingOTP, resend)

	if !resend && existingOTP != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "an unexpired OTP already exists, please check your email or request to resend the OTP",
		})
	}

	if existingOTP != "" {

		otp := existingOTP

		// Send OTP via email with existing OTP
		err = tools.SendVerificationEmail(email, otp)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to send verification email",
			})
		}

	} else {

		// No existing OTP, generate a new one and send it
		otp, err := tools.GenerateOTP()

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to generate OTP",
			})
		}

		_, err = h.DB.Exec(
			c.Context(),
			`INSERT INTO email_verification_codes (user_id, otp, expires_at) VALUES ($1, $2, NOW() + INTERVAL '60 minutes')`,
			userID,
			otp,
		)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to create OTP",
			})
		}

		err = tools.SendVerificationEmail(email, otp)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to send verification email",
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "verification email sent successfully",
	})
}
