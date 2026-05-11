package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Init(ctx context.Context) (*pgx.Conn, error) {
	msgConn := os.Getenv("DATABASE_URL")
	if msgConn == "" {
		msgConn = "postgres://postgres:postgres@localhost:5432/app_dev?sslmode=disable"
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
