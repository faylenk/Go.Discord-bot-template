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
		log.Fatal("Erro ao abrir o banco de dados:", err)
	}

	statement, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS guilds (
			id TEXT PRIMARY KEY,
			language TEXT DEFAULT 'en'
		)
	`)
	if err != nil {
		log.Fatal("Erro ao preparar a tabela:", err)
	}
	statement.Exec()
	statement.Close()

	log.Println("Banco de dados inicializado com sucesso.")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetLanguage(guildID string) string {
	if DB == nil {
		// Fallback para DEFAULT_LANG do .env
		defaultLang := os.Getenv("DEFAULT_LANG")
		if defaultLang == "" {
			return "en" // fallback final
		}
		return defaultLang
	}

	var lang string
	err := DB.QueryRow("SELECT language FROM guilds WHERE id = ?", guildID).Scan(&lang)
	if err != nil {
		// Se não encontrar guild, insere com DEFAULT_LANG
		defaultLang := os.Getenv("DEFAULT_LANG")
		if defaultLang == "" {
			defaultLang = "en"
		}
		SetLanguage(guildID, defaultLang)
		return defaultLang
	}

	return lang
}

func SetLanguage(guildID, language string) {
	if DB == nil {
		log.Println("Banco de dados não inicializado")
		return
	}

	statement, err := DB.Prepare(`
		INSERT OR REPLACE INTO guilds (id, language) VALUES (?, ?)
	`)
	if err != nil {
		log.Println("Erro ao preparar a query:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(guildID, language)
	if err != nil {
		log.Println("Erro ao salvar idioma:", err)
		return
	}
}
