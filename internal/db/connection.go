package db

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MafiaDB struct {
	db *gorm.DB
}

func NewMafiaDB(config db.MafiaBotDBConfig) (*MafiaDB, error) {
	dsn := config.DSN
	connetion, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &MafiaDB{db: connetion}, nil
}
