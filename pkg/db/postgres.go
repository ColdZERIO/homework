package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
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

func MigrationRun(ctx context.Context, db *gorm.DB) error {
	envMsg := os.Getenv("DB_MIGRATION_URL")
	if envMsg == "" {
		return fmt.Errorf("env file is empty")
	}

	m, err := migrate.New("file://pkg/migrations", envMsg)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	hasUsersTable := db.WithContext(ctx).Migrator().HasTable("users")
	if !hasUsersTable {
		version, _, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return fmt.Errorf("get migration version: %w", err)
		}
		if err == nil && version == 0 {
			if err := m.Force(-1); err != nil {
				return fmt.Errorf("force migration version: %w", err)
			}
		}
	}

	err = m.Up()
	if err == nil {
		return nil
	}

	if err != migrate.ErrNoChange {
		return fmt.Errorf("run migration: %w", err)
	}

	if hasUsersTable {
		return nil
	}

	if err := m.Force(-1); err != nil {
		return fmt.Errorf("force migration version: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rerun migration: %w", err)
	}

	return nil
}
