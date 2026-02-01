package migrate

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(dbPool *pgxpool.Pool) {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		`,
	}

	for _, q := range queries {
		_, err := dbPool.Exec(context.Background(), q)
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
	}

	log.Println("Migrations completed")
}
