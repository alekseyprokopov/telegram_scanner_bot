package storage

import (
	"scanner_bot/config"
)
type Storage interface {
	Save(p *config.Configuration) error
	Update(p *config.Configuration) error
	PickConfig(chatId int) (*config.Configuration, error)
	IsExists(chatId int) (bool, error)
}


