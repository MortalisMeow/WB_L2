package main

import (
	"log"
	"net/http"
)

func main() {
	config := LoadConfig()
	storage := NewStorage()
	server := NewServer(storage)

	handler := server.SetupRoutes()

	loggedHandler := loggingMiddleware(handler)

	log.Printf("Server starting on port %s", config.Port)
	if err := http.ListenAndServe(config.Address(), loggedHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
