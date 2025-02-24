package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rand.Seed(42) // Seed for random number generation

	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/validate", validateURLHandler)
	http.HandleFunc("/analytics/", analyticsHandler)
	http.HandleFunc("/url", getURLEntryHandler)
	http.HandleFunc("/urls", getURLsHandler)
	http.HandleFunc("/delete", deleteURLHandler)
	http.HandleFunc("/", redirectHandler) // Catch-all handler

	address := getEnv("SERVER_ADDRESS", "localhost:8080")
	log.Printf("Server is running on http://%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
