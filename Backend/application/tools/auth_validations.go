package tools

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
