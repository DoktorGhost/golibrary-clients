package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DoktorGhost/golibrary-clients/internal/repositories/postgres/dao"
	"github.com/lib/pq"
)

func (s *UsersRepository) CreateUser(user dao.UserTable) (int, error) {
	var id int
	query := `INSERT INTO users (username, password_hash, full_name) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRow(context.Background(), query, user.Username, user.PasswordHash, user.FullName).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("пользователь с таким именем уже существует")
		}
		return 0, fmt.Errorf("ошибка добавления записи: %v", err)
	}

	return id, nil
}

func (s *UsersRepository) GetUserByUsername(username string) (dao.UserTable, error) {
	var result dao.UserTable
	query := `SELECT id, username, password_hash, full_name FROM users WHERE username = $1`
	err := s.db.QueryRow(context.Background(), query, username).Scan(&result.ID, &result.Username, &result.PasswordHash, &result.FullName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dao.UserTable{}, fmt.Errorf("пользователь с username %s не найден", username)
		}
		return dao.UserTable{}, fmt.Errorf("ошибка получения пользователя: %v", err)
	}
	return result, nil
}

func (s *UsersRepository) GetUserByID(userID int) (string, error) {
	var result string
	query := `SELECT username FROM users WHERE id = $1`
	err := s.db.QueryRow(context.Background(), query, userID).Scan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("пользователь с ID %d не найден", userID)
		}
		return "", fmt.Errorf("ошибка получения пользователя: %v", err)
	}
	return result, nil
}
