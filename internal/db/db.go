package db

import (
	"log"

	"budgex_backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	return gdb, nil
}

// Simple ping helper (used in /healthz)
func Ping(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Example automigrate later:
// func AutoMigrate(gdb *gorm.DB) error {
// 	return gdb.AutoMigrate(&models.Category{}, &models.Transaction{}, &models.Budget{})
// }

func Must(gdb *gorm.DB, err error) *gorm.DB {
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	return gdb
}

func AutoMigrate(gdb *gorm.DB) error {
	// Required for gen_random_uuid()
	if err := gdb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		return err
	}
	// Tables
	if err := gdb.AutoMigrate(&models.Category{}, &models.Transaction{}, &models.Budget{}); err != nil {
		return err
	}
	// 🔑 Composite unique index for budgets upsert
	return gdb.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_budgets_user_month_category
		ON budgets (user_id, month, category_id);
	`).Error
}
