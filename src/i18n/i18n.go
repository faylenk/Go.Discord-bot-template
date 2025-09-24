package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var locales = map[string]map[string]string{}

func LoadLocales() {
	// ✅ Corrigido: caminho para a pasta locales
	files, err := filepath.Glob("src/i18n/locales/*.json")
	if err != nil {
		log.Printf("Erro ao carregar idiomas: %v", err)
		return
	}

	for _, file := range files {
		lang := filepath.Base(file)
		lang = lang[:len(lang)-5] // remove ".json"

		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Erro ao ler %s: %v", file, err)
			continue
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			log.Printf("Erro ao decodificar %s: %v", file, err)
			continue
		}

		locales[lang] = messages
		log.Printf("Idioma carregado: %s (%d mensagens)", lang, len(messages))
	}
}

func T(lang, key string) string {
	if messages, ok := locales[lang]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}
	if messages, ok := locales["en"]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}
	return key
}
