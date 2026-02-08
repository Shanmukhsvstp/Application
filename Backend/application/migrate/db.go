package migrate

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(dbPool *pgxpool.Pool) {
	queries := []string{
		`
		CREATE SEQUENCE IF NOT EXISTS user_id_seq;
		CREATE OR REPLACE FUNCTION gen_user_id()
		RETURNS BIGINT AS $$
		DECLARE
		    epoch BIGINT := 1770249600000; -- 2024-01-01 UTC
		    now_ms BIGINT;
		    seq BIGINT;
		BEGIN
		    now_ms := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT;
		    seq := nextval('user_id_seq') % 4096; -- 12 bits

		    RETURN ((now_ms - epoch) << 12) | seq;
		END;
		$$ LANGUAGE plpgsql;
		`,
		`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY DEFAULT gen_user_id(),
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			is_verified BOOLEAN NOT NULL DEFAULT FALSE,
			status TEXT NOT NULL DEFAULT 'offline',
			about TEXT DEFAULT '',
			created_at TIMESTAMP DEFAULT NOW()
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS chats (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			is_dm BOOLEAN NOT NULL DEFAULT TRUE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS messages (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			chat_id BIGINT NOT NULL,
			sender_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
			FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE NO ACTION
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS email_verification_codes (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id BIGINT NOT NULL,
			otp TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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
