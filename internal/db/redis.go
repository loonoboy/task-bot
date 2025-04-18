package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"task-bot/pkg/logger"
	"time"
)

var RDB *redis.Client
var ctx = context.Background()

// Подключение к Redis
func ConnectRedis(redisAddr string) {
	log := logger.GetLogger()
	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // адрес Redis
		Password: "",        // без пароля
		DB:       0,         // используем базу данных 0
	})

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Error("Не удалось подключиться к Redis", zap.Error(err))
	} else {
		log.Info("Подключение к Redis успешно!")
	}
}

// Закрытие соединения с Redis
func CloseRedis() {
	log := logger.GetLogger()
	if RDB != nil {
		err := RDB.Close()
		if err != nil {
			log.Error("Ошибка при закрытии подключения к Redis", zap.Error(err))
		} else {
			log.Info("Соединение с Redis закрыто.")
		}
	}
}

func CreateRedisRecord(remindAt time.Time, reminderID int64) {
	log := logger.GetLogger()
	key := fmt.Sprintf("reminder:%s:%d", remindAt.Format(time.RFC3339), reminderID)

	err := RDB.Set(ctx, key, remindAt.Format(time.RFC3339), 0).Err()
	if err != nil {
		log.Info("Ошибка записи в Redis", zap.Error(err))
	}

}

func WurmUpRedis() {
	log := logger.GetLogger()
	reminders, err := GetAllRemindersForRedis()
	if err != nil {
		log.Error("Ошибка при получении напоминаний из базы", zap.Error(err))
	}
	for _, reminder := range reminders {
		CreateRedisRecord(reminder.DueDate, reminder.ID)
	}
}
