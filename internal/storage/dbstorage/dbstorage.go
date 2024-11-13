package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	err := s.DB.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by login: %w", err)
	}

	return user, nil
}

// CardAdd метод для добавления карты.
func (s *Storage) CardAdd(ctx context.Context, card models.Card) error {
	query := `INSERT INTO cards 
	(
		user_id,
		card_bank, 
		card_number, 
		card_holder, 
		card_expiration_date, 
		card_cvv, 
		created_at
	) 
	VALUES($1,$2,$3,$4,$5,$6,now());`
	_, err := s.DB.Exec(
		ctx,
		query,
		card.UserID,
		card.CardBank,
		card.CardNumber,
		card.CardHolder,
		card.CardExpirationDate,
		card.CardCVV,
	)
	if err != nil {
		return fmt.Errorf("failed to add card: %w", err)
	}
	return nil
}

// CardGetList метод для получения списка карт.
func (s *Storage) CardGetList(ctx context.Context, userID int64) ([]models.Card, error) {
	query := `SELECT 
	card_id, 
	user_id, 
	card_bank, 
	card_number, 
	card_holder, 
	card_expiration_date, 
	card_cvv
	FROM cards WHERE user_id=$1 AND deleted_at IS NULL;`
	
	var cards []models.Card
	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card list: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error during query execution: %w", rows.Err())
	}
	defer rows.Close()
	
	for rows.Next() {
		var card models.Card
		err := rows.Scan(&card.CardID, &card.UserID, &card.CardBank, &card.CardNumber, &card.CardHolder, &card.CardExpirationDate, &card.CardCVV)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}
	return cards, nil
}

// CardDelete метод для удаления карты.
func (s *Storage) CardDelete(ctx context.Context, cardID int64) error {
	query := `UPDATE cards SET deleted_at=$1 WHERE card_id=$2;`
	result, err := s.DB.Exec(ctx, query, time.Now(), cardID)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}
	
	// Проверяем, была ли обновлена хотя бы одна запись
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("card with id %d not found", cardID)
	}
	
	return nil
}

// Close закрытие подключения к базе данных.
func (s *Storage) Close() {
	s.DB.Close()
}
