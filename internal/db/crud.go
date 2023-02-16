package db

import (
	"errors"
	"gorm.io/gorm"
)

func UpdateOrCreateUser(db MafiaGormInterface, user *User) (*User, error) {
	foundedUser := &User{}

	result := db.Where(&User{PolemicaId: user.PolemicaId}).First(foundedUser)
	if result.Error != nil && !(errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		return user, &MafiaBotUpdateOrCreateUserGetError{Detail: result.Error.Error()}
	}
	user.ID = foundedUser.ID
	user.CreatedAt = foundedUser.CreatedAt
	result = db.Save(user)
	if result.Error != nil {
		return user, &MafiaBotUpdateOrCreateUserSaveError{Detail: result.Error.Error()}
	}
	return user, nil
}
