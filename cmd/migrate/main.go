package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	envFile := flag.String("env", ".env", "Path to environment file")
	action := flag.String("action", "up", "Action: up, down, reset")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Printf("No %s file found, relying on environment variables", *envFile)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	switch *action {
	case "down":
		log.Println("Dropping tables...")
		if err := db.Migrator().DropTable(&domain.Todo{}, &domain.User{}); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		log.Println("Tables dropped.")
	case "reset":
		log.Println("Resetting database...")
		if err := db.Migrator().DropTable(&domain.Todo{}, &domain.User{}); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		log.Println("Tables dropped.")
		fallthrough
	case "up":
		log.Println("Migrating database...")
		if err := db.AutoMigrate(&domain.User{}, &domain.Todo{}); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Database migrated successfully.")
	default:
		log.Fatal("Invalid action. Use up, down, or reset.")
	}
}
