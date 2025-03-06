package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Функция для загрузки конфигурации
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Ошибка: TELEGRAM_BOT_TOKEN не задан!")
	}
	// Здесь могут быть другие настройки
}
