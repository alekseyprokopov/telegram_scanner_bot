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
	q := `INSERT INTO configs (user_id, user_config) VALUES (?, ?)`
	_, err := s.db.Exec(q, p.UserId, p.UserConfig)
	if err != nil {
		return fmt.Errorf("can's save config: %w", err)
	}
	return nil
}
func (s *Storage) Update(p *storage.Configuration) error {
	q := `UPDATE configs SET user_config = ? WHERE user_id = ?`
	_, err := s.db.Exec(q, p.UserId, p.UserConfig)
	if err != nil {
		return fmt.Errorf("can's update config: %w", err)
	}
	return nil
}

func (s *Storage) PickConfig(userId int) (*storage.Configuration, error) {
	q := `SELECT user_config FROM configs WHERE user_id = ?`

	var userConfigData string
	err := s.db.QueryRow(q, userId).Scan(&userConfigData)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("miss config: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick config: %w", err)
	}

	userConfig, err := storage.StringToConfig(userConfigData)

	return &storage.Configuration{
		UserId:     userId,
		UserConfig: *userConfig,
	}, nil

}
