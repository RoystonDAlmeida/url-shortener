package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions { // Handle preflight requests
		enableCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		fmt.Println("Received POST request")
		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "Please provide a valid URL", http.StatusBadRequest)
			return
		}

		urlShortener.mu.Lock()
		defer urlShortener.mu.Unlock() // Ensure unlocking even if an error occurs

		var existingShortURL string
		err := urlShortener.db.QueryRow("SELECT short_url FROM urls WHERE original_url = ?", originalURL).Scan(&existingShortURL)

		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if existingShortURL != "" {
			// If the original URL exists, return the existing short URL
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "The original URL already exists in the database. http://%s/%s\n", getEnv("SERVER_ADDRESS", "localhost:8080"), existingShortURL)
			return
		}

		var shortURL string
		for {
			shortURL = generateShortURL()
			var exists bool
			err := urlShortener.db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = ?)", shortURL).Scan(&exists)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if !exists { // If the short URL does not exist in the database, break out of the loop
				break
			}
		}
		fmt.Printf("Generated Short URL:%s\n", shortURL)

		stmt, err := urlShortener.db.Prepare("INSERT INTO urls(short_url, original_url) VALUES(?, ?)")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		if _, err := stmt.Exec(shortURL, originalURL); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "http://%s/%s\n", getEnv("SERVER_ADDRESS", "localhost:8080"), shortURL)
	} else {
		http.ServeFile(w, r, "index.html")
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] // Extracting the short URL from path
	fmt.Printf("%s", shortURL)

	urlShortener.mu.Lock()
	defer urlShortener.mu.Unlock()

	var originalURL string
	err := urlShortener.db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortURL).Scan(&originalURL)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

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
