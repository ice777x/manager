package utils

import (
	"database/sql"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func Config() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectDB() (*sql.DB, error) {
	dsn := "host=localhost user=root password=pass dbname=dbname port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	return db, err
}
