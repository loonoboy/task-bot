package handler

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/http"
	"task-bot/internal/db"
	"task-bot/pkg/logger"
	"time"
)

var log = logger.GetLogger()

type UserState struct {
	Step     string
	TempTask Task
}

type Task struct {
	Title       string
	Description string
	DueDate     time.Time
}

var userStates = make(map[int64]*UserState)

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

		if update.Message == nil && update.CallbackQuery == nil {
			http.Error(w, "Пустой update", http.StatusBadRequest)
			return
		}

		var userID int64
		if update.Message != nil {
			userID = update.Message.From.ID
		} else {
			userID = update.CallbackQuery.From.ID
		}

		// Если пользователя нет в userStates — создаём новую структуру
		if _, exists := userStates[userID]; !exists {
			userStates[userID] = &UserState{}
		}
		state := userStates[userID]

		if update.Message != nil {
			ProcessCommand(state, bot, update.Message)
		} else if update.CallbackQuery != nil {
			ProcessCallbackQuery(bot, update.CallbackQuery)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ProcessCommand(state *UserState, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	cmd := msg.Text
	user := msg.From
	switch {
	case cmd == "/start":
		CheckUser(user)
		response := tgbotapi.NewMessage(msg.Chat.ID,
			"Привет! 👋\n\nЯ - твой личный менеджер задач🤖\n\n"+
				"Ты можешь воспользоваться командой /help, чтобы узнать, что я умею\n\n"+
				"Выберите действие:")
		response.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "/create"),
				tgbotapi.NewInlineKeyboardButtonData("Список твоих задач", "/list"),
			),
		)
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	case cmd == "/help":
		response := tgbotapi.NewMessage(msg.Chat.ID,
			"📑 Доступные команды:\n"+
				"/create\n"+
				"/list\n")
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	case cmd == "/create" || state.Step != "":

		if userStates[user.ID].Step == "" {
			userStates[user.ID] = &UserState{Step: "waiting_for_title"}
		}
		ProcessCreate(userStates[user.ID], bot, msg)
	case cmd == "/list":
		tasks, err := db.GetUserTasks(user.ID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка получения списка задач. Попробуйте позже."))
			return
		}

		if len(tasks) == 0 {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "📭 У вас пока нет задач."))
			return
		}

		var responseText string
		for _, task := range tasks {
			responseText += fmt.Sprintf("📌 *%s*\n%s\n\n", task.Title, task.Description)
		}

		msg := tgbotapi.NewMessage(msg.Chat.ID, responseText)
		msg.ParseMode = "Markdown"
		bot.Send(msg)

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
	switch callback.Data {
	case "/create":
		CheckUser(callback.From)
		if userStates[callback.From.ID].Step == "" {
			userStates[callback.From.ID] = &UserState{Step: "waiting_for_title"}
		}
		ProcessCreate(userStates[callback.From.ID], bot, callback.Message)
	case "/list":
		tasks, err := db.GetUserTasks(callback.From.ID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(callback.Message.Chat.ID, "❌ Ошибка получения списка задач. Попробуйте позже."))
			return
		}

		if len(tasks) == 0 {
			bot.Send(tgbotapi.NewMessage(callback.Message.Chat.ID, "📭 У вас пока нет задач."))
			return
		}

		var responseText string
		for _, task := range tasks {
			responseText += fmt.Sprintf(
				"📌 Название: %s \nОписание: %s\n🕒 Время: %s\n",
				task.Title,
				task.Description,
				task.DueDate.Format("02.01.2006 15:04"),
			)
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, responseText)
		bot.Send(msg)
	}
}

func ProcessCreate(state *UserState, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch state.Step {
	case "waiting_for_title":
		state.Step = "waiting_for_description"
		response := tgbotapi.NewMessage(msg.Chat.ID, "📝 Введите название задачи:")
		_, err := bot.Send(response)
		if err != nil {
			log.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	case "waiting_for_description":
		state.TempTask.Title = msg.Text
		state.Step = "waiting_for_date"
		response := tgbotapi.NewMessage(msg.Chat.ID, "📌 Теперь введите описание задачи:")
		bot.Send(response)

	case "waiting_for_date":
		state.TempTask.Description = msg.Text
		state.Step = "waiting_for_setup"
		response := tgbotapi.NewMessage(msg.Chat.ID, "📥 Пожалуйста, введите дату и время в формате:\n"+
			"ДД.ММ.ГГГГ ЧЧ:ММ\n"+
			"Например: 10.04.2025 14:30")
		bot.Send(response)
	case "waiting_for_setup":
		layout := "02.01.2006 15:04"
		loc, _ := time.LoadLocation("Europe/Moscow")
		parsedTime, err := time.ParseInLocation(layout, msg.Text, loc)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Неправильый формат, введите еще раз"))
			return
		} else if parsedTime.Before(time.Now()) {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Напоминание должно быть на будущее"))
			return
		}
		reminderID, err := db.CreateTask(msg.From.ID, state.TempTask.Title, state.TempTask.Description, parsedTime)
		db.CreateRedisRecord(parsedTime, reminderID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка сохранения задачи. Попробуйте позже."))
			return
		}
		delete(userStates, msg.From.ID)
		response := tgbotapi.NewMessage(msg.Chat.ID, "✅ Задача успешно создана!")
		bot.Send(response)
	}
}

func CheckUser(user *tgbotapi.User) {
	if d, err := db.CheckUserExistence(user.ID); err != nil {
		log.Error("Ошибка при проверке существования пользователя", zap.Error(err))
	} else if !d {
		if err := db.AddUser(user.ID, user.FirstName, user.LastName, user.UserName); err != nil {
			log.Error("Ошибка создания пользователя", zap.Error(err))
		}
	} else {
		log.Info("Пользователь уже сущесвтует")
	}
}
