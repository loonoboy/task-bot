-- init.sql

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL ,
                                     telegram_id BIGINT PRIMARY KEY,
                                     first_name VARCHAR(255),
                                     last_name VARCHAR(255),
                                     username VARCHAR(255),
                                     language_code VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS tasks (
                                    id SERIAL,
                                    user_id BIGINT REFERENCES users(telegram_id) ON DELETE CASCADE PRIMARY KEY,
                                    title VARCHAR(255) NOT NULL,
                                    description TEXT,
                                    status_send BOOLEAN DEFAULT false,
                                    due_date TIMESTAMP
);