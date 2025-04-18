package db

import (
	"context"
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

type Task struct {
	ID          int64
	UserID      int
	Title       string
	Description string
	DueDate     time.Time
}

var DB *pgxpool.Pool
var log = logger.GetLogger()

func ConnectDB(dsn string) {
	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Error("Не удалось подключиться к БД", zap.Error(err))
	}
	log.Info("Подключение к PostgreSQL успешно!")

}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Info("Подключение к БД закрыто.")
	}
}

func AddUser(telegramID int64, firstName, lastName, username string) error {
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

func CreateTask(userID int64, title, description string, duedate time.Time) (int64, error) {
	var reminderID int64
	log := logger.GetLogger()
	query := `INSERT INTO tasks (user_id, title, description, due_date) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	err := DB.QueryRow(context.Background(), query, userID, title, description, duedate).Scan(&reminderID)
	if err != nil {
		log.Error("failed to insert task", zap.Error(err))
		return reminderID, err
	}
	log.Info("Task created successfully")
	return reminderID, nil
}

func GetUserTasks(userID int64) ([]Task, error) {
	ctx := context.Background()
	rows, err := DB.Query(ctx, `
        SELECT title, description, due_date
        FROM tasks 
        WHERE user_id = $1 and status_send = false`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.Title, &t.Description, &t.DueDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTaskByID(taskID int64) (*Task, error) {
	var task Task
	ctx := context.Background()
	err := DB.QueryRow(ctx, `
        SELECT user_id, title, description
        FROM tasks 
        WHERE id = $1`, taskID).Scan(&task.UserID, &task.Title, &task.Description)
	if err != nil {
		log.Error("failed to insert task", zap.Error(err))
		return nil, err
	}
	return &task, nil
}

func UpdateStatusSend(reminderID int64) error {
	ctx := context.Background()
	query := `UPDATE tasks SET status_send = $1 WHERE id = $2`
	_, err := DB.Exec(ctx, query, true, reminderID)
	if err != nil {
		log.Error("failed to insert task", zap.Error(err))
		return err
	}
	return nil
}

func GetAllRemindersForRedis() ([]Task, error) {
	ctx := context.Background()
	rows, err := DB.Query(ctx, `
        SELECT id, due_date
        FROM tasks WHERE status_send = false`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.DueDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
