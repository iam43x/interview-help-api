package store

import (
	"fmt"
	"log/slog"

	"github.com/iam43x/interview-help-api/internal/domain"
)

func (s *Store) ExistsInvite(value string) error {
	i := &domain.Invite{}
	query := "SELECT value, isActive FROM invites WHERE value=?;"
	row := s.db.QueryRow(query, value)
	if err := row.Scan(&i.Value, &i.IsActive); err != nil {
		return fmt.Errorf("ошибка поиска инвайта %w", err)
	}
	if !i.IsActive {
		return fmt.Errorf("Invites deactivate value=%v", value)
	}
	return nil
}

func (s *Store) DeactivateInvite(value string) {
	query := "UPDATE invites SET isActive=false WHERE value=?;"
	res, err := s.db.Exec(query, value)
	if err != nil {
        slog.Error("error while deactivate invite [%v]: %w ", value, err)
    }
    rowsAffected, err := res.RowsAffected()
    if err != nil || rowsAffected != 1 {
        slog.Error("error while deactivate invite [%v]: %w ", value, err)
    }
    
    slog.Debug("Updated invite value=%v.", value)
}