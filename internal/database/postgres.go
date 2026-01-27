package database

import (
	"fmt"
	"log"
	"time"

	"github.com/datmedevil17/go-vuln/internal/config"
	"github.com/datmedevil17/go-vuln/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ PostgreSQL connected successfully")

	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("⚠️  Warning: Failed to create UUID extension: %v", err)
	}

	// In development mode, drop and recreate tables to ensure schema is correct
	// This fixes issues with column name changes (e.g., git_hub_id -> github_id)
	// In development mode, auto-migration will handle schema updates.
	// Previously we dropped tables here, but that causes data loss on restart.

	// Auto-migrate database schema
	log.Println("🔄 Running auto-migration...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Repository{},
		&models.Workflow{},
		&models.ScanResult{},
		&models.WorkflowExecution{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	log.Println("✅ Database schema migrated successfully")

	return db, nil
}
