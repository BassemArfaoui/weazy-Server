package config

import (
	"fmt"
	"os"

	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	var err error
	godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connected!")
	// Migrate(DB)
	return DB, nil
}


func Migrate(db *gorm.DB) error {

	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	err := db.AutoMigrate(
		&models.Chat{},
		
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	fmt.Println("Database migrated successfully!")
	return nil
}
