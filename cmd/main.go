package main

import (
	"go.uber.org/zap"
	"net/http"
	"task-bot/internal/bot"
	"task-bot/internal/db"
	"task-bot/internal/router"
	"task-bot/pkg/config"
	"task-bot/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	db.ConnectDB()
	defer db.CloseDB()
	logger.InitLogger()
	log := logger.GetLogger()
	defer func() {
		if err := log.Sync(); err != nil {
			log.Error("Ошибка при сбросе логов", zap.Error(err))
		}
	}()
	tgBot, err := bot.NewBot(cfg.BotToken, cfg.WebhookURL)
	if err != nil {
		log.Error("Ошибка при запуске бота:", zap.Error(err))
	}
	tgBot.API.Debug = cfg.Debug
	bot.SetBotMenu(tgBot.API)

	r := router.SetupRouter(tgBot.API)

	addr := ":8443"
	log.Info("Запуск сервера", zap.String("address", addr))

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Ошибка сервера", zap.Error(err))
	}
}
