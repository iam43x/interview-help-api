package store

import (
	"fmt"

	"github.com/iam43x/interview-help-api/internal/domain"
)

func (s *Store) GetUserByLogin(login string) (*domain.User, error) {
	u := &domain.User{}
	query := "SELECT login, name, pass FROM users WHERE login=?;"
	row := s.db.QueryRow(query, login)
	if err := row.Scan(&u.Login, &u.Name, &u.Pass); err != nil {
		return nil, fmt.Errorf("ошибка поиска пользователя %w", err)
	}
	return u, nil
}

func (s *Store) CreateUser(login, name, pass, invite string) (*domain.User, error) {
	query := "INSERT INTO users (login, name, pass, invite) VALUES (?, ?, ?, ?);"
	_, err := s.db.Exec(query, login, name, pass, invite)
	if err != nil {
		return nil, fmt.Errorf("ошибка вставки в таблицу: %w", err)
	}
	return &domain.User{
		Login: login,
		Name: name,
	}, nil
}