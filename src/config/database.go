package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./bot.db"
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Error at opening DB:", err)
	}

	statement, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS guilds (
			id TEXT PRIMARY KEY,
			language TEXT DEFAULT 'en'
		)
	`)
	if err != nil {
		log.Fatal("Error at making table:", err)
	}
	statement.Exec()
	statement.Close()

	log.Println("Database initialized successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func EnsureGuildExists(guildID string) {
	if DB == nil {
		return
	}

	// Verify if the guild already exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM guilds WHERE id = ?", guildID).Scan(&count)
	if err != nil {
		log.Println("Error checking guild existence:", err)
		return
	}

	if count == 0 {
		SetLanguage(guildID, getDefaultLang())
		log.Printf("Created guild entry for %s with default language\n", guildID)
	}
}

func GetLanguage(guildID string) string {
	if DB == nil {
		return getDefaultLang()
	}

	// If guildID is empty, return default language
	if guildID == "" {
		return getDefaultLang()
	}

	var lang string
	err := DB.QueryRow("SELECT language FROM guilds WHERE id = ?", guildID).Scan(&lang)

	if err != nil {
		// If no language is set, use the default language
		defaultLang := getDefaultLang()
		SetLanguage(guildID, defaultLang)
		return defaultLang
	}

	return lang
}

func getDefaultLang() string {
	defaultLang := os.Getenv("DEFAULT_LANG")
	if defaultLang == "" {
		return "en"
	}
	return defaultLang
}

func SetLanguage(guildID, language string) {
	if DB == nil {
		log.Println("Database not initialized. Cannot set language.")
		return
	}

	statement, err := DB.Prepare(`
		INSERT OR REPLACE INTO guilds (id, language) VALUES (?, ?)
	`)
	if err != nil {
		log.Println("Error at setting query:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(guildID, language)
	if err != nil {
		log.Println("Error saving language:", err)
		return
	}
}
