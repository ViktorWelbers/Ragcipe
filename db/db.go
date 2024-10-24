package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	// Connection string for postgres
	connStr := "host=localhost port=5432 user=postgres password=your_secure_password dbname=mydatabase sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL!")
	return db, nil
}
