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

	// 2. Build the connection string (DSN)
	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("--- DSN is empty! Using hardcoded local fallback...")
		// Use the Key-Value format (replace 'your_password' with your actual 4-dot password)
		dsn = "host=localhost user=postgres password=1111 dbname=expense-tracker port=5432 sslmode=disable"
	}

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
