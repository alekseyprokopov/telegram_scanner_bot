package storage

import (
	"scanner_bot/config"
)
type Storage interface {
	Save(p *config.Configuration) error
	Update(chatId int64, UserConfig string) error
	GetConfig(chatId int64) (*config.Configuration, error)
	IsExists(chatId int64) (bool, error)
}


