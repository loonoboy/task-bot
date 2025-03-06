package router

import (
	"github.com/go-chi/chi/v5"
	"task-bot/task-bot/internal/handler"
)

// Создаём новый роутер
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handler.DefoultHandler)
	r.Get("/start", handler.StartHandler) // Пример маршрута для команды start
	return r
}
