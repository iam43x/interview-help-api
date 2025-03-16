package store

import (
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(filepath string) (*Store, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}
	log.Printf("sqlite3 %s started!\n", filepath)
	return &Store{db: db}, err
}