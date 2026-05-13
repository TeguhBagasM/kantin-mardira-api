package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	if DB != nil {
		return DB
	}

	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS revoked_tokens (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			jti VARCHAR(64) UNIQUE NOT NULL,
			user_id UUID NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			transaction_code VARCHAR(50) UNIQUE NOT NULL,
			customer_name VARCHAR(100),
			cashier_id UUID REFERENCES users(id),
			payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('cash', 'qris')),
			payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
			total_amount INTEGER NOT NULL CHECK (total_amount >= 0),
			paid_amount INTEGER DEFAULT 0,
			change_amount INTEGER DEFAULT 0,
			transaction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		ALTER TABLE IF EXISTS transactions
		ADD COLUMN IF NOT EXISTS customer_name VARCHAR(100)
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_items (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
			menu_id UUID REFERENCES menus(id),
			quantity INTEGER NOT NULL CHECK (quantity > 0),
			price INTEGER NOT NULL CHECK (price >= 0),
			subtotal INTEGER NOT NULL CHECK (subtotal >= 0)
		)
	`).Error; err != nil {
		log.Fatal(err)
	}

	DB = db
	log.Println("Database Connected")
	return DB
}