package handler

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/http"
	"task-bot/pkg/logger"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func WebHookHandler(bot *tgbotapi.BotAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
			return
		}

		if update.Message != nil {
			ProcessMessage(bot, update.Message)
		}
		if update.CallbackQuery != nil {
			ProcessCallbackQuery(bot, update.CallbackQuery)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ProcessMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log := logger.GetLogger()
	cmd := msg.Command()
	switch cmd {
	case "start":
		response := tgbotapi.NewMessage(msg.Chat.ID,
			"Привет! 👋\n\nЯ - твой личный менеджер задач🤖\n\n"+
				"Ты можешь воспользоваться командой /help, чтобы узнать, что я умею\n\n"+
				"Выберите действие:")
		response.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("start", "start"),
				tgbotapi.NewInlineKeyboardButtonData("help", "help"),
			),
		)
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	case "help":
		response := tgbotapi.NewMessage(msg.Chat.ID,
			"📑 Доступные команды:")
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	default:
		response := tgbotapi.NewMessage(msg.Chat.ID,
			"⛔ Неизвестная команда\n\nВоспользуйся командой /help, чтобы узнать какие команды я могу выполнять")
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	}

}

func ProcessCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	log := logger.GetLogger()
	switch callback.Data {
	case "start":
		response := tgbotapi.NewMessage(callback.Message.Chat.ID,
			"Привет!👋 Я твой личный менеджер задач. Помогу запомнить важные события и напомню о них")
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	case "help":
		response := tgbotapi.NewMessage(callback.Message.Chat.ID,
			"Привет! Команды которые я знаю:")
		response.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("start", "start"),
				tgbotapi.NewInlineKeyboardButtonData("help", "help"),
			),
		)
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения: %v", zap.Error(err))
		}
	}
}
