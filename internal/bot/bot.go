package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"task-bot/pkg/logger"
)

// Bot структура для хранения объекта бота
type Bot struct {
	API *tgbotapi.BotAPI
}

// NewBot создаёт и настраивает нового Telegram-бота
func NewBot(BotToken, webhookURL string) (*Bot, error) {
	log := logger.GetLogger()
	botAPI, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Info("Ошибка создания Webhook", zap.Error(err))
	}

	_, err = botAPI.Request(webhookConfig)
	if err != nil {
		log.Fatal("Ошибка установки Webhook", zap.Error(err))
	}

	log.Info("Webhook установлен:", zap.String("webhook_url", webhookURL))
	log.Info("Авторизован", zap.String("user_name", botAPI.Self.UserName))

	return &Bot{API: botAPI}, nil
}

func SetBotMenu(bot *tgbotapi.BotAPI) {
	log := logger.GetLogger()
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Справка"},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(cfg)
	if err != nil {
		log.Error("Ошибка установки меню", zap.Error(err))
	}
}
