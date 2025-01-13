package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbHost     = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
)

const (
	TableMessages = "messages"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Connected to PostgreSQL!")
	DB = db

	// Initialize tables
	err = initializeTables(db)
	if err != nil {
		return nil, fmt.Errorf("error initializing tables: %v", err)
	}

	return db, nil
}

func initializeTables(db *sql.DB) error {
	// Create messages table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id SERIAL PRIMARY KEY,
            content TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
            host_ip TEXT NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
