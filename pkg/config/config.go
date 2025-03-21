package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
	"task-bot/pkg/logger"
)

// Config — структура для хранения конфигурации
type Config struct {
	BotToken   string
	Debug      bool
	WebhookURL string
}

// LoadConfig загружает конфигурацию из .env
func LoadConfig() *Config {
	log := logger.GetLogger()
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Info("Предупреждение: файл .env не найден, загружаем переменные из окружения")
	}

	// Читаем переменные
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	debugStr := os.Getenv("DEBUG")
	webhookURL := os.Getenv("WEBHOOK")

	// Преобразуем DEBUG в bool
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		log.Error("Ошибка при разборе DEBUG, установлено значение false", zap.Error(err))
		debug = false
	}

	return &Config{
		BotToken:   botToken,
		Debug:      debug,
		WebhookURL: webhookURL,
	}
}
