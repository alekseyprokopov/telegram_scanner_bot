package sqllite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"scanner_bot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Save(p *storage.Configuration) error {
	q := `INSERT INTO configs (username, user_config) VALUES (?, ?)`
	_, err := s.db.Exec(q, p.UserName, p.UserConfig)
	if err != nil {
		return fmt.Errorf("can's save config: %w", err)
	}
	return nil
}

func (s *Storage) Pick(userName string) (*storage.Configuration, error) {
	q := `SELECT user_config FROM configs WHERE username = ?`
	result, err := s.db.Exec(q, userName)
	if err != nil {
		return nil, fmt.Errorf("can't pick config")
	}

}
