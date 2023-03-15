package db

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal"
	"github.com/google/uuid"
)

type MafiaBotDBError struct {
	internal.MafiaBotError
}

type MafiaBotToGameResultError struct {
	MafiaBotDBError
	WinnerCode int
}

func (e *MafiaBotToGameResultError) Error() string {
	return fmt.Sprintf("%v: Not supported winner code: %d", e.GetISOFormat(), e.WinnerCode)
}

type MafiaBotUpdateOrCreateUserGetError struct {
	MafiaBotDBError
	Detail string
}

func (e *MafiaBotUpdateOrCreateUserGetError) Error() string {
	return fmt.Sprintf("%v: Error while get user from db: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotUpdateOrCreateUserSaveError struct {
	MafiaBotDBError
	Detail string
}

func (e *MafiaBotUpdateOrCreateUserSaveError) Error() string {
	return fmt.Sprintf("%v: Error while save user object to db: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotGetLastGameDriverError struct {
	MafiaBotDBError
	Detail string
}

func (e *MafiaBotGetLastGameDriverError) Error() string {
	return fmt.Sprintf("%v: Error while retrieve last game: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotGetLastGameEmptyDBError struct {
	MafiaBotDBError
}

func (e *MafiaBotGetLastGameEmptyDBError) Error() string {
	return fmt.Sprintf("%v: Empty game table", e.GetISOFormat())
}

type MafiaBotGetDailyStatisticQueryError struct {
	MafiaBotDBError
	Detail string
}

func (e *MafiaBotGetDailyStatisticQueryError) Error() string {
	return fmt.Sprintf("%v: Can't execute the query: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotGetDailyStatisticRowScanError struct {
	MafiaBotDBError
	Detail string
}

func (e *MafiaBotGetDailyStatisticRowScanError) Error() string {
	return fmt.Sprintf("%v: Can't read row: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotMarkGamesAsSentFindError struct {
	MafiaBotDBError
	Games  []uuid.UUID
	Detail string
}

func (e *MafiaBotMarkGamesAsSentFindError) Error() string {
	return fmt.Sprintf("%v: Can't get info about games %s: %s", e.GetISOFormat(), e.Games, e.Detail)
}

type MafiaBotMarkGamesAsSentTransactionError struct {
	MafiaBotDBError
	Games  []uuid.UUID
	Detail string
}

func (e *MafiaBotMarkGamesAsSentTransactionError) Error() string {
	return fmt.Sprintf("%v: Can't mark games %s as sent: %s", e.GetISOFormat(), e.Games, e.Detail)
}
