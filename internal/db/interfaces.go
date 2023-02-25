package db

import (
	"database/sql"
	"gorm.io/gorm"
)

type MafiaDBInterface interface {
	UpdateOrCreateUser(user *User) (*User, error)
	GetLastGame() (*Game, error)
	Create(value interface{}) error
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}
