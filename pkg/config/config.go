package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config — структура для хранения конфигурации
type Config struct {
	BotToken   string
	Debug      bool
	WebhookURL string
}

// LoadConfig загружает конфигурацию из .env
func LoadConfig() *Config {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: файл .env не найден, загружаем переменные из окружения")
	}

	// Читаем переменные
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	debugStr := os.Getenv("DEBUG")
	webhookURL := os.Getenv("WEBHOOK")

	// Преобразуем DEBUG в bool
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		log.Println("Ошибка при разборе DEBUG, установлено значение false")
		debug = false
	}

	return &Config{
		BotToken:   botToken,
		Debug:      debug,
		WebhookURL: webhookURL,
	}
}
