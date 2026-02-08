package tools

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/resend/resend-go/v3"
)

func UsernameIsUnique(username string, db *pgxpool.Pool) (bool, error) {

	count := 0

	err := db.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users WHERE username=$1",
		username,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil
	}

	return true, nil
}

func UserAlreadyExist(email string, db *pgxpool.Pool) (bool, error) {

	count := 0

	err := db.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users WHERE email=$1",
		email,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func SendVerificationEmail(email string, otp string) error {
	// Implement email sending logic here using an email service provider
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil // or return an error indicating that the API key is missing
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		To:      []string{email},
		From:    "no-reply@aureliahealth.org",
		Subject: "Application's Email Verification",
		Html: fmt.Sprintf(
			`<p>OTP for email verification: <strong>%s</strong></p>`,
			otp,
		),
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	return err
}

func GenerateOTP() (string, error) {

	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
