package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"scanner_bot/config"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	//if err := db.Ping(); err != nil {
	//	return nil, fmt.Errorf("can't connect to database: %w", err)
	//}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Save(p *config.Configuration) error {
	q := `INSERT INTO configs (chat_id, user_config) VALUES (?, ?)`
	_, err := s.db.Exec(q, p.ChatId, p.UserConfig)
	if err != nil {
		return fmt.Errorf("can's save config: %w", err)
	}
	return nil
}
func (s *Storage) Update(p *config.Configuration) error {
	q := `UPDATE configs SET user_config = ? WHERE chat_id = ?`
	_, err := s.db.Exec(q, p.ChatId, p.UserConfig)
	if err != nil {
		return fmt.Errorf("can's update config: %w", err)
	}
	return nil
}

func (s *Storage) GetConfig(chatId int) (*config.Configuration, error) {
	q := `SELECT user_config FROM configs WHERE chat_id = ?`

	var userConfigData string
	err := s.db.QueryRow(q, chatId).Scan(&userConfigData)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("miss config: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick config: %w", err)
	}

	userConfig, err := config.StringToConfig(userConfigData)

	return &config.Configuration{
		ChatId:     chatId,
		UserConfig: *userConfig,
	}, nil
}

func (s *Storage) IsExists(chatId int) (bool, error) {
	q := `SELECT COUNT(*) FROM configs WHERE chat_id = ?`

	var count int

	if err := s.db.QueryRow(q, chatId).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}
	return count > 0, nil
}

func (s *Storage) Init() error {
	query := `CREATE TABLE IF NOT EXISTS configs (user_name TEXT, user_config TEXT)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}