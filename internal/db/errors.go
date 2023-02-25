package db

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal"
)

type MafiaBotEnumError struct {
	internal.MafiaBotError
}

type MafiaBotToGameResultError struct {
	MafiaBotEnumError
	WinnerCode int
}

func (e *MafiaBotToGameResultError) Error() string {
	return fmt.Sprintf("%v: Not supported winner code: %d", e.GetISOFormat(), e.WinnerCode)
}

type MafiaBotUpdateOrCreateUserGetError struct {
	MafiaBotEnumError
	Detail string
}

func (e *MafiaBotUpdateOrCreateUserGetError) Error() string {
	return fmt.Sprintf("%v: Error while get user from db: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotUpdateOrCreateUserSaveError struct {
	MafiaBotEnumError
	Detail string
}

func (e *MafiaBotUpdateOrCreateUserSaveError) Error() string {
	return fmt.Sprintf("%v: Error while save user object to db: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotGetLastGameDriverError struct {
	MafiaBotEnumError
	Detail string
}

func (e *MafiaBotGetLastGameDriverError) Error() string {
	return fmt.Sprintf("%v: Error while retrieve last game: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotGetLastGameEmptyDBError struct {
	MafiaBotEnumError
}

func (e *MafiaBotGetLastGameEmptyDBError) Error() string {
	return fmt.Sprintf("%v: Empty game table", e.GetISOFormat())
}
