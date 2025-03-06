package main

import (
	"log"
	"task-bot/task-bot/internal/bot"
	"task-bot/task-bot/pkg/config"
)

func main() {
	// Загружаем конфигурацию
	config.Load()

	// Создаём и запускаем бота
	tgBot, err := bot.NewBot()
	if err != nil {
		log.Fatal("Ошибка запуска бота:", err)
	}

	tgBot.Start()
}
