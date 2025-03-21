package router

import (
	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"task-bot/internal/handler"
)

// Создаём новый роутер
func SetupRouter(bot *tgbotapi.BotAPI) *chi.Mux {
	r := chi.NewRouter()
	r.Use(handler.LoggingMiddleware)

	r.Post("/webhook", handler.WebHookHandler(bot))
	return r
}
