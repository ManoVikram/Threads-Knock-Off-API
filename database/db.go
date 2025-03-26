package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connectionString := os.Getenv("DATABASE_URL")

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %s", err))
	}

	if err = DB.Ping(); err != nil {
		panic(fmt.Sprintf("Database ping failed: %s", err))
	}

	log.Println("Connected to database successfully")
}