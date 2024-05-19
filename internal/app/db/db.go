package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func GetPostgresConnectionString() string {
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_db := os.Getenv("POSTGRES_DB")
	connStr := "user=%s password=%s dbname=%s sslmode=disable"
	connStr = fmt.Sprintf(connStr, postgres_user, postgres_password, postgres_db)
	return connStr
}

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_db := os.Getenv("POSTGRES_DB")

	connStr := "user=%s password=%s dbname=%s sslmode=disable"
	connStr = fmt.Sprintf(connStr, postgres_user, postgres_password, postgres_db)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db, nil
}
