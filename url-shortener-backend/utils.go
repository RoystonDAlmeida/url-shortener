package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allows all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func generateShortURL() string {
	charset := getEnv("CHARSET", "abcdefghijklmnopqrstuvwxyz")
	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortURL)
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Function to log click events
func logClickEvent(shortURL string, r *http.Request) {

	ipAddress := r.RemoteAddr  // Get IP address of the requester
	userAgent := r.UserAgent() // Get user agent string

	// Use urlShortener.db to access the database
	_, err := urlShortener.db.Exec("INSERT INTO clicks (short_url, ip_address, user_agent) VALUES (?, ?, ?)", shortURL, ipAddress, userAgent)
	if err != nil {
		log.Printf("Error logging click event: %v", err)
	}
}

// Helper function to format date from string to desired format.
func formatDate(expirationDate string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, expirationDate)
	if err != nil {
		return "", err
	}
	return parsedTime.Format("02 Jan 2006"), nil
}
