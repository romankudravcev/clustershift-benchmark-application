package main

import (
	"benchmarker/api"
	"benchmarker/db"
	"log"
)

func main() {
	// Connect to DB
	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer client.Close()

	router := api.SetupRouter()
	if err := router.Run(":1337"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
