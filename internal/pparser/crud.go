package pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"gorm.io/gorm"
)

type ParserMafiaDB struct {
	db.MafiaDB
}

func (pdb *ParserMafiaDB) SaveMinimalGameStatistic(gameStatistic MinimalGameStatistic) error {
	err := pdb.MafiaDB.Db.Transaction(func(tx *gorm.DB) error {
		players := [10]*db.User{}
		for idx, p := range gameStatistic.Players {
			user, err := db.UpdateOrCreateUser(pdb.MafiaDB.Db, &db.User{PolemicaId: p.ID, PolemicaNickName: p.NickName})
			if err != nil {
				return &MafiaBotPolemicaParserSaveMinimalGameStatisticUserError{Detail: err.Error()}
			}
			players[idx] = user
		}
		game := &db.Game{
			PolemicaId: gameStatistic.ID,
			StartedAt:  gameStatistic.StartedAt,
			Winner:     gameStatistic.GameResult,
		}
		result := pdb.Db.Create(game)
		if result.Error != nil {
			return &MafiaBotPolemicaParserSaveMinimalGameStatisticGameError{Detail: result.Error.Error()}
		}

		for idx, p := range gameStatistic.Players {
			playerGame := &db.PlayerGame{UserID: players[idx].ID, GameID: game.ID, Points: p.Score}
			result = pdb.Db.Create(playerGame)
			if result.Error != nil {
				return &MafiaBotPolemicaParserSaveMinimalGameStatisticPlayerGameError{Detail: result.Error.Error()}
			}
		}
		return nil
	})
	return err
}
