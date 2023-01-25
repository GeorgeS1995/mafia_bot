package pparser

import "github.com/GeorgeS1995/mafia_bot/internal/db"

type ParserMafiaDB struct {
	db.MafiaDB
}

func (d *ParserMafiaDB) SaveMinimalGameStatistic(gameStatistic MinimalGameStatistic) error {
	return nil
}
