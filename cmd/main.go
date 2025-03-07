package main

import (
	"log"
	"net/http"
	"task-bot/internal/bot"
	"task-bot/internal/router"
	"task-bot/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	// Создаёмs и запускаем бота
	tgBot, err := bot.NewBot(cfg.BotToken, cfg.WebhookURL)
	if err != nil {
		log.Fatal("Ошибка при запуске бота:", err)
	}
	tgBot.API.Debug = cfg.Debug

	r := router.SetupRouter(tgBot.API)

	log.Println("Сервер запущен на порту 8443...")
	log.Println("Бот успешно запущен")
	log.Fatal(http.ListenAndServe(":8443", r))
}
