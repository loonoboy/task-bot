package main

import (
	"go.uber.org/zap"
	"net/http"
	"task-bot/internal/bot"
	"task-bot/internal/router"
	"task-bot/pkg/config"
	"task-bot/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger()
	log := logger.GetLogger()
	defer log.Sync()
	// Создаёмs и запускаем бота
	tgBot, err := bot.NewBot(cfg.BotToken, cfg.WebhookURL)
	if err != nil {
		log.Error("Ошибка при запуске бота:", zap.Error(err))
	}
	tgBot.API.Debug = cfg.Debug

	r := router.SetupRouter(tgBot.API)

	addr := ":8443"
	log.Info("Запуск сервера", zap.String("address", addr))

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Ошибка сервера", zap.Error(err))
	}
}
