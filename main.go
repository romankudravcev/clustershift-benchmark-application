package main

import (
	"benchmarker/api"
	"benchmarker/db"
	"log"
)

func main() {
	// Connect to DB
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Defer close for postgres only
	if db.DB != nil {
		defer db.DB.Close()
	}

	router := api.SetupRouter()
	if err := router.Run(":1337"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
