package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) *pgxpool.Pool {

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Database connected")

	return pool
}
