package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *sql.DB  // Para SQL puro (clientes)
	GormDB *gorm.DB // Para GORM (produtos)
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetConfig returns database configuration from environment variables
func GetConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "gestao_user"),
		Password: getEnv("DB_PASSWORD", "gestao_pass"),
		DBName:   getEnv("DB_NAME", "gestao_db"),
	}
}

// Connect establishes database connections (both SQL and GORM)
func Connect() error {
	config := GetConfig()
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	
	// Conectar SQL puro (para clientes)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Conectar GORM (para produtos)
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect GORM: %w", err)
	}
	
	log.Printf("Successfully connected to database: %s (SQL + GORM)", config.DBName)
	return nil
}

// Close closes both database connections
func Close() error {
	if DB != nil {
		if err := DB.Close(); err != nil {
			return err
		}
	}
	
	// GORM connection will be closed automatically
	// when the underlying sql.DB is closed
	if GormDB != nil {
		sqlDB, err := GormDB.DB()
		if err == nil {
			return sqlDB.Close()
		}
	}
	
	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}