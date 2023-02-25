package db

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

func (b *MafiaDB) UpdateOrCreateUser(user *User) (*User, error) {
	foundedUser := &User{}

	result := b.Db.Where(&User{PolemicaId: user.PolemicaId}).First(foundedUser)
	if result.Error != nil && !(errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		return user, &MafiaBotUpdateOrCreateUserGetError{Detail: result.Error.Error()}
	}
	user.ID = foundedUser.ID
	user.CreatedAt = foundedUser.CreatedAt
	result = b.Db.Save(user)
	if result.Error != nil {
		return user, &MafiaBotUpdateOrCreateUserSaveError{Detail: result.Error.Error()}
	}
	return user, nil
}

func (b *MafiaDB) GetLastGame() (*Game, error) {
	game := &Game{}
	result := b.Db.Order("started_at desc").First(game)
	var err error
	if result.Error != nil && !(errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		err = &MafiaBotGetLastGameDriverError{Detail: result.Error.Error()}
	} else if result.Error != nil && (errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		err = &MafiaBotGetLastGameEmptyDBError{}
	}
	return game, err
}

func (b *MafiaDB) Create(value interface{}) error {
	result := b.Db.Create(value)
	return result.Error
}

func (b *MafiaDB) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return b.Db.Transaction(fc, opts...)
}
