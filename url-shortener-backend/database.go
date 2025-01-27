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
	createTable()
}

func createTable() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS urls (
	    short_url TEXT PRIMARY KEY,
	    original_url TEXT NOT NULL
	    );`
	if _, err := urlShortener.db.Exec(sqlStmt); err != nil {
		log.Fatal(err)
	}
}
