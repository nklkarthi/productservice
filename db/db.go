package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"productservice/models"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	migrateSchema()
}

func migrateSchema() {
	err := DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("Error migrating schema: %v", err)
	}
}
