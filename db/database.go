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

		// Configure MongoDB connection pool options
		clientOptions := options.Client().ApplyURI(mongoURI)
		clientOptions.SetMaxPoolSize(25)                         // Maximum number of connections
		clientOptions.SetMinPoolSize(5)                          // Minimum number of connections
		clientOptions.SetMaxConnIdleTime(1 * time.Minute)        // Maximum idle time
		clientOptions.SetConnectTimeout(10 * time.Second)        // Connection timeout
		clientOptions.SetServerSelectionTimeout(5 * time.Second) // Server selection timeout

		client, err := mongo.Connect(ctx, clientOptions)
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

	// Configure connection pool to prevent connection exhaustion
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime
	db.SetConnMaxIdleTime(1 * time.Minute) // Maximum idle time

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
            host_ip INET NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	// Add indexes for better query performance
	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at DESC);
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_messages_host_ip ON messages(host_ip);
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
