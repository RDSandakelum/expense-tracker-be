package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. Get database connection credentials from Environment Variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Fallback to local defaults if env variables are empty
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "1111"
	}
	if dbName == "" {
		dbName = "expense-tracker"
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	// 2. Build the connection string (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Colombo",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// 3. Open the connection
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	// 4. Automatically sync your Go structs with database tables (Auto Migration)
	err = database.AutoMigrate(
		&User{},
		&Account{},
		&Category{},
		&SubCategory{},
		&Budget{},
		&Goal{},
		&Transaction{},
		&SavingsWithdrawal{},
		&AccountTransfer{},
		&BudgetTemplate{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}

	// 5. Assign to the global DB variable
	DB = database
}
