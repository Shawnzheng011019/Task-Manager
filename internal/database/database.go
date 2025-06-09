package database

import (
	"fmt"
	"task-manager/internal/config"
	"task-manager/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	cfg := config.GetConfig()
	
	var err error
	switch cfg.Database.Type {
	case "sqlite3":
		DB, err = gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := DB.AutoMigrate(&models.Task{}, &models.SubTask{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
