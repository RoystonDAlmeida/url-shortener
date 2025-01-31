package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type URLEntry struct {
	ShortURL       string  `json:"short_url"`
	LongURL        string  `json:"long_url"`
	CustomAlias    *string `json:"custom_alias,omitempty"` // Pointer to allow null values
	ExpirationDate *string `json:"expiration_date,omitempty"`
}

// URL structure to hold each url information for getURLshandler()
type URL struct {
	ShortURL       string `json:"short_url"`
	OriginalURL    string `json:"original_url"`
	CustomAlias    string `json:"custom_alias"`
	ExpirationDate string `json:"expiration_date"`
	TotalClicks    int    `json:"total_clicks"`
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for all requests
	enableCORS(w)
	if r.Method == http.MethodOptions { // Handle preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		fmt.Println("Received POST request")
		originalURL := r.FormValue("url")
		customAlias := r.FormValue("alias")         // Get custom alias from form data
		expirationDate := r.FormValue("expiration") // Get expiration date from form data

		if originalURL == "" {
			http.Error(w, "Please provide a valid URL", http.StatusBadRequest)
			return
		}

		urlShortener.mu.Lock()
		defer urlShortener.mu.Unlock() // Ensure unlocking even if an error occurs

		var existingShortURL string
		var existingCustomAlias string
		var existingExpirationDate string

		fmt.Println("Checking if original URL exists in the database...")
		err := urlShortener.db.QueryRow("SELECT short_url, custom_alias, expiration_date FROM urls WHERE original_url = ?", originalURL).Scan(&existingShortURL, &existingCustomAlias, &existingExpirationDate)

		if err != nil && err != sql.ErrNoRows {
			log.Printf("Database query error: %v", err) // Log database errors
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		fmt.Println("Checked original URL existence.")

		if existingShortURL != "" {
			if expirationDate != "" {
				// If form expirationDate is provided
				existingExpirationTime, err := time.Parse("2006-01-02T00:00:00Z", existingExpirationDate)

				if err != nil {
					http.Error(w, "Invalid expiration date format", http.StatusBadRequest)
					return
				}

				now := time.Now()

				if existingExpirationTime.Before(now) || existingExpirationTime.Truncate(24*time.Hour).Equal(now.Truncate(24*time.Hour)) {
					stmt, err := urlShortener.db.Prepare("UPDATE urls SET expiration_date = ? WHERE original_url = ?")
					if err != nil {
						log.Printf("Error preparing update statement: %v", err) // Log preparation errors
						http.Error(w, "Database error", http.StatusInternalServerError)
						return
					}
					defer stmt.Close()

					if _, err := stmt.Exec(expirationDate, originalURL); err != nil {
						log.Printf("Error executing update statement: %v", err) // Log execution errors
						http.Error(w, "Database error", http.StatusInternalServerError)
						return
					}

					w.Header().Set("Access-Control-Allow-Origin", "*")
					formattedExpirationDate, _ := formatDate(expirationDate)

					fmt.Fprintf(w, "The expiration date for the existing URL has been updated to %s. Short URL: http://%s/%s\n",
						formattedExpirationDate, getEnv("SERVER_ADDRESS", "localhost:8080"), existingShortURL)
					return
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					fmt.Fprintf(w, "The provided expiration date is not valid for updating. Current expiration date remains unchanged. Short URL: http://%s/%s\n",
						getEnv("SERVER_ADDRESS", "localhost:8080"), existingShortURL)
					return
				}
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "The original URL already exists in the database. Short URL: http://%s/%s\n",
				getEnv("SERVER_ADDRESS", "localhost:8080"), existingShortURL)
			return
		}

		var shortURL string
		maxAttempts := 10
		attempts := 0

		if customAlias != "" {
			var exists bool
			err := urlShortener.db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE custom_alias = ?)", customAlias).Scan(&exists)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "Custom alias already in use", http.StatusBadRequest)
				return
			}
			shortURL = customAlias // Use the provided custom alias
		} else {
			for attempts < maxAttempts {
				shortURL = generateShortURL()
				var exists bool
				err := urlShortener.db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = ?)", shortURL).Scan(&exists)
				if err != nil {
					log.Printf("Database query error: %v", err) // Log database errors
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}

				if !exists { // If the short URL does not exist in the database, break out of the loop
					break
				}
				attempts++
			}

			if attempts == maxAttempts {
				http.Error(w, "Unable to generate a unique short URL after multiple attempts", http.StatusInternalServerError)
				return
			}
		}

		fmt.Printf("Generated Short URL: %s\n", shortURL)

		var expiration *string = nil
		if expirationDate != "" {
			expiration = &expirationDate
		}

		stmt, err := urlShortener.db.Prepare("INSERT INTO urls(short_url, original_url, custom_alias, expiration_date) VALUES(?, ?, ?, ?)")
		if err != nil {
			log.Printf("Error preparing statement: %v", err) // Log preparation errors
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		if _, err := stmt.Exec(shortURL, originalURL, customAlias, expiration); err != nil {
			log.Printf("Error executing statement: %v", err) // Log execution errors
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "New Short URL created: http://%s/%s\n",
			getEnv("SERVER_ADDRESS", "localhost:8080"), shortURL)
	} else {
		http.ServeFile(w, r, "index.html")
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortURL := r.URL.Path[1:] // Extracting the short URL from path
	fmt.Printf("%s", shortURL)

	urlShortener.mu.Lock()
	defer urlShortener.mu.Unlock()

	var originalURL string
	var expirationDate sql.NullString // Use sql.NullString to handle nullable fields

	err := urlShortener.db.QueryRow("SELECT original_url, expiration_date FROM urls WHERE short_url = ?", shortURL).Scan(&originalURL, &expirationDate)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if the link has expired
	if expirationDate.Valid {
		expirationTime, err := time.Parse(time.RFC3339, expirationDate.String) // Parse the expiration date string
		if err != nil {
			http.Error(w, "Invalid expiration date format", http.StatusInternalServerError)
			return
		}
		if time.Now().After(expirationTime) { // Check if current time is after expiration time
			http.Error(w, "This link has expired", http.StatusGone)
			return
		}
	}

	// Log the click event
	logClickEvent(shortURL, r)

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func validateURLHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Enable CORS for this handler

	if r.Method == http.MethodOptions {
		return // Respond to preflight requests
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	urlToValidate := r.FormValue("url")
	resp, err := http.Get(urlToValidate)
	if err != nil {
		http.Error(w, "URL is not reachable", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Valid URL"))
	} else {
		http.Error(w, "URL is not reachable", http.StatusBadRequest)
	}
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Enable CORS for this handler

	if r.Method == http.MethodOptions {
		return // Respond to preflight requests
	}

	shortURL := r.URL.Path[len("/analytics/"):]

	urlShortener.mu.Lock()
	defer urlShortener.mu.Unlock()

	data, err := getAnalyticsData(shortURL)
	if err != nil {
		http.Error(w, "Error retrieving analytics data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data) // Send JSON response with analytics data
}

func getURLEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for all requests
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		shortUrl := r.URL.Query().Get("short_url") // Extract short_url from query parameters

		if shortUrl == "" {
			http.Error(w, "Missing short_url parameter", http.StatusBadRequest)
			return
		}

		//
		var urlEntry URLEntry

		// Query the database for the URL entry
		err := urlShortener.db.QueryRow("SELECT short_url, original_url, custom_alias, expiration_date FROM urls WHERE short_url = ?", shortUrl).
			Scan(&urlEntry.ShortURL, &urlEntry.LongURL, &urlEntry.CustomAlias, &urlEntry.ExpirationDate)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No entry found for the provided short URL", http.StatusNotFound)
				return
			}
			log.Printf("Database query error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Format expiration date if it exists
		if urlEntry.ExpirationDate != nil {
			parsedDate, err := time.Parse("2006-01-02T15:04:05Z", *urlEntry.ExpirationDate) // Parse the incoming date
			if err != nil {
				log.Printf("Date parsing error: %v", err)
				http.Error(w, "Invalid expiration date format", http.StatusBadRequest)
				return
			}
			formattedDate := parsedDate.Format("02 January 2006") // Format to desired output
			urlEntry.ExpirationDate = &formattedDate              // Update expiration date in the struct
		}

		// Set response header and encode the response as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urlEntry)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getURLsHandler handles the GET request to fetch all URLs
func getURLsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Enable CORS for pre-flight requests

	rows, err := urlShortener.db.Query("SELECT short_url, original_url, custom_alias, expiration_date, (SELECT COUNT(*) FROM clicks WHERE clicks.short_url = urls.short_url) AS total_clicks FROM urls")
	if err != nil {
		http.Error(w, "Failed to query URLs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.ShortURL, &url.OriginalURL, &url.CustomAlias, &url.ExpirationDate, &url.TotalClicks); err != nil {
			http.Error(w, "Failed to scan URLs", http.StatusInternalServerError)
			return
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(urls)
}

// deleteURLHandler handles the DELETE request to delete a URL by its short URL
func deleteURLHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Enable CORS for pre-flight requests

	shortURL := r.URL.Query().Get("short_url") // Get the short URL from query parameters
	if shortURL == "" {
		http.Error(w, "Missing short URL", http.StatusBadRequest)
		return
	}

	// Check if the URL exists in the database before attempting to delete
	var originalURL string
	err := urlShortener.db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortURL).Scan(&originalURL)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows were found, return a 404 Not Found error
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		// If there was another error querying the database, return an internal server error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Proceed to delete the URL since it exists
	stmt, err := urlShortener.db.Prepare("DELETE FROM urls WHERE short_url = ?")
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(shortURL)
	if err != nil {
		http.Error(w, "Failed to execute statement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content on successful deletion
}
