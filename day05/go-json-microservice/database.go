package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "user=webuser password=webpass123 dbname=productdb sslmode=disable"
	//connStr := "user=webuser dbname=personal_info_db sslmode=disable password=webpass123 host=localhost"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Create table if not exists
	query := `
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name TEXT,
        price NUMERIC(10,2),
        in_stock BOOLEAN
    );`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Table creation failed: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL and ensured table exists.")
}
