package db

import (
	"database/sql/driver"
)

type GameResult string

const (
	Draw     GameResult = "Draw"
	CityWin  GameResult = "CityWin"
	MafiaWin GameResult = "MafiaWin"
)

func (gr *GameResult) Scan(value interface{}) error {
	*gr = GameResult(value.(string))
	return nil
}

func (gr GameResult) Value() (driver.Value, error) {
	return string(gr), nil
}

func ToGameResult(winnerCode int) (GameResult, error) {
	winnerCodeMap := map[int]GameResult{
		0: Draw,
		1: CityWin,
		2: MafiaWin,
	}
	if r, ok := winnerCodeMap[winnerCode]; !ok {
		return Draw, &MafiaBotToGameResultError{WinnerCode: winnerCode}
	} else {
		return r, nil
	}
}
