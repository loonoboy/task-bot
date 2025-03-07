package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// Bot структура для хранения объекта бота
type Bot struct {
	API *tgbotapi.BotAPI
}

// NewBot создаёт и настраивает нового Telegram-бота
func NewBot(BotToken, webhookURL string) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatalf("Ошибка создания Webhook: %v", err)
	}

	_, err = botAPI.Request(webhookConfig)
	if err != nil {
		log.Fatalf("Ошибка установки Webhook: %v", err)
	}

	log.Println("Webhook установлен:", webhookURL)
	log.Printf("Авторизован как %s", botAPI.Self.UserName)

	return &Bot{API: botAPI}, nil
}
