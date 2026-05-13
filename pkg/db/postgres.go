package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func Init(ctx context.Context) (*pgx.Conn, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		return nil, fmt.Errorf("load env file: %w", err)
// 	}

// 	msgConn := os.Getenv("DB_URL")
// 	if msgConn == "" {
// 		return nil, fmt.Errorf("env file is empty")
// 	}

// 	conn, err := pgx.Connect(ctx, msgConn)
// 	if err != nil {
// 		return nil, fmt.Errorf("connect to postgres: %w", err)
// 	}

// 	err = conn.Ping(ctx)
// 	if err != nil {
// 		conn.Close(ctx)
// 		return nil, fmt.Errorf("ping postgres: %w", err)
// 	}

// 	sqlFile, err := os.ReadFile("pkg/db/schema.sql")
// 	if err != nil {
// 		conn.Close(ctx)
// 		return nil, fmt.Errorf("read schema.sql: %w", err)
// 	}

// 	_, err = conn.Exec(ctx, string(sqlFile))
// 	if err != nil {
// 		conn.Close(ctx)
// 		return nil, fmt.Errorf("apply schema: %w", err)
// 	}

// 	return conn, nil
// }

func Init(ctx context.Context) (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	msgConn := os.Getenv("DB_URL")
	if msgConn == "" {
		return nil, fmt.Errorf("env file is empty")
	}

	conn, err := gorm.Open(postgres.Open(msgConn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	db, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("get db from gorm: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	sqlFile, err := os.ReadFile("pkg/db/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("read schema.sql: %w", err)
	}

	err = conn.Exec(string(sqlFile)).Error
	if err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return conn, nil
}
