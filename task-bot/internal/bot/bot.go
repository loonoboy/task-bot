package bot

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot структура для хранения объекта бота
type Bot struct {
	API *tgbotapi.BotAPI
}

// NewBot создаёт и настраивает нового Telegram-бота
func NewBot() (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, err
	}

	botAPI.Debug = true // Включаем режим отладки

	log.Printf("Авторизован как %s", botAPI.Self.UserName)

	return &Bot{API: botAPI}, nil
}

// Start запускает обработку сообщений
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот!")
			b.API.Send(msg)
		}
	}
}
