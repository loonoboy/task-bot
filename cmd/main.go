package main

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"task-bot/internal/bot"
	"task-bot/internal/db"
	"task-bot/internal/router"
	"task-bot/internal/service"
	"task-bot/pkg/config"
	"task-bot/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	db.ConnectDB(cfg.DSN)
	defer db.CloseDB()

	db.ConnectRedis(cfg.RedisAddr)
	defer db.CloseRedis()
	db.WurmUpRedis()

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
	service.StartReminderWatcher(context.Background(), tgBot.API, db.RDB)

	addr := ":8443"
	log.Info("Запуск сервера", zap.String("address", addr))

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Ошибка сервера", zap.Error(err))
	}
}
