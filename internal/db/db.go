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

// internal/db/db.go
func AutoMigrate(gdb *gorm.DB) error {
	if err := gdb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		return err
	}
	if err := gdb.AutoMigrate(&models.Category{}, &models.Transaction{}, &models.Budget{}); err != nil {
		return err
	}
	// 🔧 ensure user_id is TEXT in all tables
	if err := gdb.Exec(`ALTER TABLE categories   ALTER COLUMN user_id TYPE text USING user_id::text;`).Error; err != nil {
		return err
	}
	if err := gdb.Exec(`ALTER TABLE transactions ALTER COLUMN user_id TYPE text USING user_id::text;`).Error; err != nil {
		return err
	}
	if err := gdb.Exec(`ALTER TABLE budgets      ALTER COLUMN user_id TYPE text USING user_id::text;`).Error; err != nil {
		return err
	}

	if err := gdb.Exec(`
  -- transactions.category_id: text -> uuid (keep NULLs safe)
  ALTER TABLE transactions
  ALTER COLUMN category_id DROP NOT NULL;
`).Error; err != nil {
		return err
	}

	if err := gdb.Exec(`
  ALTER TABLE transactions
  ALTER COLUMN category_id TYPE uuid
  USING (CASE WHEN category_id IS NULL OR category_id = '' THEN NULL ELSE category_id::uuid END);
`).Error; err != nil {
		return err
	}

	if err := gdb.Exec(`
  -- budgets.category_id: text -> uuid
  ALTER TABLE budgets
  ALTER COLUMN category_id TYPE uuid
  USING (CASE WHEN category_id IS NULL OR category_id = '' THEN NULL ELSE category_id::uuid END);
`).Error; err != nil {
		return err
	}

	// (optional) helpful indexes for analytics
	_ = gdb.Exec(`CREATE INDEX IF NOT EXISTS idx_tx_user_date ON transactions (user_id, date);`).Error
	_ = gdb.Exec(`CREATE INDEX IF NOT EXISTS idx_tx_user_type_date ON transactions (user_id, type, date);`).Error
	_ = gdb.Exec(`CREATE INDEX IF NOT EXISTS idx_tx_user_cat_date ON transactions (user_id, category_id, date);`).Error

	// existing unique index for budgets stays valid
	return gdb.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_budgets_user_month_category
        ON budgets (user_id, month, category_id);
    `).Error
}
