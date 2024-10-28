package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage структура для бд.
type Storage struct {
	DB *pgxpool.Pool
}

// GetPingDB проверка соединения с базой данных.
func (s *Storage) GetPingDB(ctx context.Context) error {
	err := s.DB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

// RegisterUser метод принимает логин и пароль, проверяет на уникальность логин,
// сохранят в таблице users, и возвращает uid и ошибку.
func (s *Storage) RegisterUser(ctx context.Context, login, hashedPassword string) (int64, error) {
	//готовим запрос на вставку
	query := `INSERT INTO users (username, password, created_at) VALUES($1,$2,now());`
	_, err := s.DB.Exec(ctx, query, login, hashedPassword)
	//обработка ошибки сохранения нового пользователя
	if err != nil {
		//если login неуникальный
		//if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		//БД возвращает ошибку на "русском" языке. Из-за этого не обрабатывается ошибка. Как исправить не нашел.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", login, models.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: не удалось выполнить запись в базу %w", login, err)
	}

	//извлечение ID
	queryID := `SELECT user_id FROM users WHERE username=$1;`
	var uid int64
	err = s.DB.QueryRow(ctx, queryID, login).Scan(&uid)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", login, err)
	}
	if err != nil {
		return 0, fmt.Errorf("unable to execute queryID: %w", err)
	}

	return uid, nil
}

// GetUserByLogin метод для получения пользователя по логину.
func (s *Storage) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	query := `SELECT user_id, username, password, created_at FROM users WHERE username=$1;`
	var user models.User
	err := s.DB.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by login: %w", err)
	}

	return user, nil
}

// Close закрытие подключения к базе данных.
func (s *Storage) Close() {
	s.DB.Close()
}
