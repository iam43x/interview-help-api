package store

import (
	"errors"
	"fmt"

	"github.com/iam43x/interview-help-api/internal/domain"
)

var ErrInviteIsNotActive = errors.New("инвайт не активен!")

func (s *Store) ExistsInvite(value string) error {
	i := &domain.Invite{}
	query := "SELECT value, isActive FROM invites WHERE value=?;"
	row := s.db.QueryRow(query, value)
	if err := row.Scan(&i.Value, &i.IsActive); err != nil {
		return fmt.Errorf("ошибка поиска инвайта %w", err)
	}
	if !i.IsActive {
		return ErrInviteIsNotActive
	}
	return nil
}