package handler

import (
	"fmt"
	"net/http"
)

// Пример обработчика для команды start
func DefoultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hellow, i'm ur Telegram-bot!")
}
func StartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bot started")
}
