package main

import (
	"log"
	"time"
)

func getAnalyticsData(shortURL string) (map[string]interface{}, error) {
	var clickCount int

	// Query to count total clicks for the short URL
	err := urlShortener.db.QueryRow("SELECT COUNT(*) FROM clicks WHERE short_url = ?", shortURL).Scan(&clickCount)
	if err != nil {
		log.Printf("Error counting clicks: %v", err)
		return nil, err
	}

	// Query to get timestamps of all clicks for the short URL
	rows, err := urlShortener.db.Query("SELECT clicked_at FROM clicks WHERE short_url = ?", shortURL)
	if err != nil {
		log.Printf("Error retrieving click timestamps: %v", err)
		return nil, err
	}
	defer rows.Close()

	clicksByDate := make(map[string]map[string]interface{})

	for rows.Next() {
		var clickedAt string
		if err := rows.Scan(&clickedAt); err != nil {
			log.Printf("Error scanning timestamp: %v", err)
			return nil, err
		}

		// Parse the UTC timestamp (assuming it's in RFC3339 format with a 'Z' suffix)
		utcTime, err := time.Parse(time.RFC3339, clickedAt)
		if err != nil {
			log.Printf("Error parsing timestamp: %v", err)
			return nil, err
		}

		// Convert UTC time to local time based on system settings
		localTime := utcTime.In(time.Local)

		dateKey := localTime.Format("02 Jan 2006") // Format date as "29 Jan 2025"

		if _, exists := clicksByDate[dateKey]; !exists {
			clicksByDate[dateKey] = map[string]interface{}{
				"click_counts": 0,
				"timestamps":   []string{},
			}
		}

		clicksByDate[dateKey]["click_counts"] = clicksByDate[dateKey]["click_counts"].(int) + 1
		clicksByDate[dateKey]["timestamps"] = append(clicksByDate[dateKey]["timestamps"].([]string), localTime.Format(time.RFC3339))
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return map[string]interface{}{
		"total_click_counts": clickCount,
		"day":                clicksByDate,
	}, nil
}
