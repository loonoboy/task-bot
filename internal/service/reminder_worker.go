package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"task-bot/internal/db"
	"task-bot/pkg/logger"
	"time"
)

var log = logger.GetLogger()

func processReminders(bot *tgbotapi.BotAPI, rdb *redis.Client) {
	now := time.Now()
	keys, err := rdb.Keys(context.Background(), "reminder:*").Result()
	if err != nil {
		log.Info("Ошибка получения ключей:", zap.Error(err))
		return
	}

	for _, fullKey := range keys {
		key := strings.TrimPrefix(fullKey, "reminder:")

		lastColon := strings.LastIndex(key, ":")
		if lastColon == -1 {
			log.Warn("Неверный формат ключа:", zap.String("key", fullKey))
			continue
		}

		timestampStr := key[:lastColon]
		reminderIDStr := key[lastColon+1:]
		ts, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Warn("Ошибка парсинга времени:", zap.String("timestamp", timestampStr), zap.Error(err))
			continue
		}
		log.Info("", zap.Time("", now), zap.Time("", ts), zap.Bool("", now.After(ts)))
		if now.After(ts) {
			reminderID, err := strconv.Atoi(reminderIDStr)
			if err != nil {
				log.Warn("Неверный ID напоминания:", zap.String("id", reminderIDStr))
				continue
			}

			reminder, err := db.GetTaskByID(int64(reminderID))
			if err != nil {
				log.Info("Не найдено напоминание", zap.Error(err))
				continue
			}

			message := fmt.Sprintf("🔔 Напоминание!\n\nЗадача: %s\nОписание: %s", reminder.Title, reminder.Description)
			_, err = bot.Send(tgbotapi.NewMessage(int64(reminder.UserID), message))
			if err != nil {
				log.Info("Ошибка отправки сообщения", zap.Error(err))
				continue
			}

			err = db.UpdateStatusSend(int64(reminderID))
			if err != nil {
				log.Info("Ошибка обновления status_send", zap.Error(err))
			}

			// ВАЖНО: удаляем оригинальный ключ с префиксом!
			rdb.Del(context.Background(), fullKey)
		}
	}
}

func StartReminderWatcher(ctx context.Context, bot *tgbotapi.BotAPI, rdb *redis.Client) {
	log.Info("start processing reminders")
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				processReminders(bot, rdb)
			}
		}
	}()
}
