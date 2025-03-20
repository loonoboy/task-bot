package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

// InitLogger инициализирует глобальный логгер
func InitLogger() {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

// GetLogger возвращает глобальный логгер
func GetLogger() *zap.Logger {
	if log == nil {
		InitLogger()
	}
	return log
}
