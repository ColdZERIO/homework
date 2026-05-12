package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Init(ctx context.Context) (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	msgConn := os.Getenv("DB_URL")
	if msgConn == "" {
		return nil, fmt.Errorf("env file is empty")
	}

	conn, err := pgx.Connect(ctx, msgConn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	sqlFile, err := os.ReadFile("pkg/db/schema.sql")
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("read schema.sql: %w", err)
	}

	_, err = conn.Exec(ctx, string(sqlFile))
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return conn, nil
}
