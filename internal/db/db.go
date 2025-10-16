package db

import (
	"log"

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
