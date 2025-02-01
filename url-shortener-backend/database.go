package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

type URLShortener struct {
	db *sql.DB
	mu sync.Mutex // To handle concurrent access
}

var urlShortener *URLShortener

func init() {
	var err error
	urlShortener = &URLShortener{}
	dbPath := getEnv("DATABASE_URL", "../urls.db") // Get database path from env variable
	urlShortener.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Enable foreign key support
	_, err = urlShortener.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

// createTables creates the necessary tables for the URL shortener application
func createTables() {
	// SQL statement to create urls table
	sqlStmtURLs := `
	CREATE TABLE IF NOT EXISTS urls (
	    short_url TEXT PRIMARY KEY,
	    original_url TEXT NOT NULL,
	    custom_alias TEXT,
	    expiration_date DATETIME
	);`
	if _, err := urlShortener.db.Exec(sqlStmtURLs); err != nil {
		log.Fatal(err)
	}

	// SQL statement to create clicks table
	sqlStmtClicks := `
	CREATE TABLE IF NOT EXISTS clicks (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    short_url TEXT NOT NULL,
	    clicked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    ip_address TEXT,
	    user_agent TEXT,
	    FOREIGN KEY (short_url) REFERENCES urls (short_url) ON DELETE CASCADE
	);`
	if _, err := urlShortener.db.Exec(sqlStmtClicks); err != nil {
		log.Fatal(err)
	}
}
