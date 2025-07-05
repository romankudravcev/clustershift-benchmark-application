package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbHost     = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
	dbType     = os.Getenv("DB_TYPE")
	mongoURI   = os.Getenv("MONGODB_URI")
)

const (
	TableMessages = "messages"
)

var (
	DB      *sql.DB
	MongoDB *mongo.Database
)

func ConnectDB() (interface{}, *mongo.Client, error) {
	if dbType == "mongodb" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			return nil, nil, err
		}
		db := client.Database(dbName)
		MongoDB = db
		log.Println("Connected to MongoDB!")
		return db, client, nil
	}

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Connected to PostgreSQL!")
	DB = db

	// Initialize tables
	err = initializeTables(db)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing tables: %v", err)
	}

	return db, nil, nil
}

func initializeTables(db *sql.DB) error {
	// Create messages table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id SERIAL PRIMARY KEY,
            content TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
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
