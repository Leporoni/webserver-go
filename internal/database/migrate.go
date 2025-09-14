package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	// Get absolute path to migrations directory
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("could not get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// CreateMigration creates a new migration file
func CreateMigration(name string) error {
	// This would be used for creating new migration files
	// For now, we'll create them manually
	log.Printf("To create a migration, add files to migrations/ directory:")
	log.Printf("  migrations/001_%s.up.sql", name)
	log.Printf("  migrations/001_%s.down.sql", name)
	return nil
}