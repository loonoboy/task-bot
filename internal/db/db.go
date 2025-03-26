package db

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"task-bot/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID         int
	TelegramID int64
	FirstName  string
	LastName   string
	Username   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var DB *pgxpool.Pool

func ConnectDB() {
	log := logger.GetLogger()
	dsn := "postgres://taskbot:1909@localhost:5432/taskbotdb"

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Error("Не удалось подключиться к БД", zap.Error(err))
	}
	log.Info("Подключение к PostgreSQL успешно!")

}

func CloseDB() {
	log := logger.GetLogger()
	if DB != nil {
		DB.Close()
		log.Info("Подключение к БД закрыто.")
	}
}

func AddUser(telegramID int64, firstName, lastName, username string) error {
	log := logger.GetLogger()
	query := `INSERT INTO users (telegram_id, first_name, last_name, username) 
              VALUES ($1, $2, $3, $4)`
	_, err := DB.Exec(context.Background(), query, telegramID, firstName, lastName, username)
	if err != nil {
		log.Error("failed to insert user", zap.Error(err))
	} else {
		log.Info("create user")
	}
	return nil
}

func CheckUserExistence(userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE telegram_id = $1)`
	err := DB.QueryRow(context.Background(), query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateTask(userID int64, title, description string) error {
	log := logger.GetLogger()
	query := `INSERT INTO tasks (user_id, title, description) VALUES ($1, $2, $3)`
	_, err := DB.Exec(context.Background(), query, userID, title, description)
	if err != nil {
		log.Error("failed to insert task", zap.Error(err))
	}

	fmt.Println("Task created successfully")
	return nil
}
