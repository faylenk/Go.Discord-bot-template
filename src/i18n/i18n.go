package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

// `defaultLang` is used as a fallback if a translation is missing in the requested language.
const defaultLang = "en"

var locales = map[string]map[string]string{}

// LoadLocales loads all locale JSON files from the locales directory.
func LoadLocales() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("It was not possible to get the current filename")
	}
	localesPath := filepath.Join(filepath.Dir(filename), "locales", "*.json")

	files, err := filepath.Glob(localesPath)
	if err != nil {
		log.Printf("Error at searching for the language file: %v", err)
		return err
	}

	if len(files) == 0 {
		log.Printf("No language file found: %s", localesPath)
		return nil
	}

	for _, file := range files {
		lang := filepath.Base(file)
		lang = lang[:len(lang)-5] // removes ".json"

		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Erro ao ler %s: %v", file, err)
			continue
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			log.Printf("Error at decoding: %s: %v", file, err)
			continue
		}

		locales[lang] = messages
		log.Printf("Langugages loaded: %s (%d messages)", lang, len(messages))
	}

	log.Printf("Total of languages loaded: %d", len(locales))
	return nil
}

// T translates a message key into the specified language.
func T(lang, key string) string {
	if messages, ok := locales[lang]; ok {
		if message, ok := messages[key]; ok {
			return message
		}
	}

	if lang != defaultLang {
		if messages, ok := locales[defaultLang]; ok {
			if message, ok := messages[key]; ok {
				return message
			}
		}
	}
	return key
}

// GetAvailableLocales returns a sorted list of available locale codes.
func GetAvailableLocales() []string {
	keys := make([]string, 0, len(locales))
	for k := range locales {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Order the keys for consistency
	return keys
}
