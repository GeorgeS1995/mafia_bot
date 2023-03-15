package db

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MafiaDBInterface interface {
	UpdateOrCreateUser(user *User) (*User, error)
	GetLastGame() (*Game, error)
	Create(value interface{}) error
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	GetDailyStatistic(DaySwitchHour int) ([]*DailyStatistic, error)
	MarkGamesAsSent(gameIDS []uuid.UUID) (err error)
}
