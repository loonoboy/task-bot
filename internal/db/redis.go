package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"task-bot/pkg/logger"
)

var redisClient *redis.Client
var ctx = context.Background()

// Подключение к Redis
func ConnectRedis() {
	log := logger.GetLogger()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // адрес Redis
		Password: "",               // без пароля
		DB:       0,                // используем базу данных 0
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Error("Не удалось подключиться к Redis", zap.Error(err))
	} else {
		log.Info("Подключение к Redis успешно!")
	}
}

// Закрытие соединения с Redis
func CloseRedis() {
	log := logger.GetLogger()
	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			log.Error("Ошибка при закрытии подключения к Redis", zap.Error(err))
		} else {
			log.Info("Соединение с Redis закрыто.")
		}
	}
}
