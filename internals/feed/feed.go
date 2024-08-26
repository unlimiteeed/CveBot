package feed

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

const DiscordWebhookURL = ""

func stripHTMLTags(input string) string {
	// Unescape HTML entities
	unescaped := html.UnescapeString(input)

	// Regular expression to match HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleanedString := re.ReplaceAllString(unescaped, "")

	return cleanedString
}

func initDB(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS cve (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cve_id TEXT UNIQUE,
		title TEXT,
		published TEXT,
		description TEXT,
		severity TEXT,
		link TEXT
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}
}

func cveExists(db *sql.DB, cveID string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM cve WHERE cve_id=? LIMIT 1)"
	err := db.QueryRow(query, cveID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Error checking if CVE exists: %s", err)
	}
	return exists
}

func insertCVE(db *sql.DB, cveID, title, published, description, severity, link string) {
	stmt, err := db.Prepare("INSERT INTO cve(cve_id, title, published, description, severity, link) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Error preparing insert statement: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(cveID, title, published, description, severity, link)
	if err != nil {
		log.Fatalf("Error inserting data: %s", err)
	}
}

func sendDiscordNotification(cveID, title, published, description, severity, link string) {
	embed := map[string]interface{}{
		"title":       fmt.Sprintf("New CVE Alert: %s", cveID),
		"description": title,
		"url":         link,
		"color":       0xff0000, // Red color for severity
		"fields": []map[string]interface{}{
			{
				"name":   "Published",
				"value":  published,
				"inline": true,
			},
			{
				"name":   "Severity",
				"value":  severity,
				"inline": true,
			},
			{
				"name":  "Description",
				"value": description,
			},
		},
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error creating JSON payload: %s", err)
	}

	resp, err := http.Post(DiscordWebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error sending Discord notification: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Fatalf("Unexpected response from Discord: %d", resp.StatusCode)
	}

	fmt.Println("Notification sent to Discord for CVE:", cveID)
}

func Readfunc() {

	db, err := sql.Open("sqlite3", "./cve_data.db")
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()

	// Initialize the database schema
	initDB(db)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://cvefeed.io/rssfeed/latest.xml")
	if err != nil {
		fmt.Println("Error parsing feed:", err)
		return
	}

	// Determine the number of items to process (up to 10)
	numItems := 10
	if len(feed.Items) < 10 {
		numItems = len(feed.Items)
	}

	// Loop through the latest 10 items and process them
	for i := 0; i < numItems; i++ {
		item := feed.Items[i]
		cveID := item.Title // Assuming the CVE ID is part of the title
		title := item.Title
		published := item.Published
		description := stripHTMLTags(item.Description)
		severity := item.Custom["severity"] // Assuming severity is stored in Custom map
		link := item.Link

		// Check if the CVE already exists in the database
		if !cveExists(db, cveID) {
			// Send a Discord notification with embed
			sendDiscordNotification(cveID, title, published, description, severity, link)

			// Insert the new CVE into the database
			insertCVE(db, cveID, title, published, description, severity, link)
		} else {
			fmt.Println("CVE already exists in the database:", cveID)
		}
	}

	fmt.Println("Finished processing CVEs.")
}
