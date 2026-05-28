package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(ctx context.Context) (*gorm.DB, error) {
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

	return conn, nil
}

func MigrationRun() error {
	envMsg := os.Getenv("DB_MIGRATION_URL")
	if envMsg == "" {
		return fmt.Errorf("env file is empty")
	}

	m, err := migrate.New("file://pkg/migrations", envMsg)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rerun migration: %w", err)
	}

	return nil
}
