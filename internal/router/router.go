package router

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/http"
	"task-bot/pkg/logger"
)

// Создаём новый роутер
func SetupRouter(bot *tgbotapi.BotAPI) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
			return
		}

		if update.Message != nil {
			ProcessMessage(bot, update.Message)
		}
		w.WriteHeader(http.StatusOK)
	})
	return r
}

func ProcessMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log := logger.GetLogger()
	switch msg.Text {
	case "/start":
		response := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Я твой бот.")
		_, err := bot.Send(response)
		if err != nil {
			log.Info("Ошибка отправки сообщения: %v", zap.Error(err))
		}
	case "/help":
		response := tgbotapi.NewMessage(msg.Chat.ID, "/start, /help")
		_, err := bot.Send(response)
		if err != nil {
			log.Info("Ошибка отправки сообщения: %v", zap.Error(err))
		}
	case "Hello":
		response := tgbotapi.NewMessage(msg.Chat.ID, "Hello, "+msg.From.FirstName)
		_, err := bot.Send(response)
		if err != nil {
			log.Info("Ошибка отправки сообщения: %v", zap.Error(err))
		}
	default:
		log.Info("Неизвестная команда: %s", zap.String("massage_text", msg.Text))
	}
}
