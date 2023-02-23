package pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"gorm.io/gorm"
)

type MafiaParserServiceHandler struct {
	db.MafiaDBInterface
}

func (pdb *MafiaParserServiceHandler) SaveMinimalGameStatistic(gameStatistic MinimalGameStatistic) error {
	err := pdb.Transaction(func(tx *gorm.DB) error {
		players := [10]*db.User{}
		for idx, p := range gameStatistic.Players {
			user, err := pdb.UpdateOrCreateUser(&db.User{PolemicaId: p.ID, PolemicaNickName: p.NickName})
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
		err := pdb.Create(game)
		if err != nil {
			return &MafiaBotPolemicaParserSaveMinimalGameStatisticGameError{Detail: err.Error()}
		}

		for idx, p := range gameStatistic.Players {
			playerGame := &db.PlayerGame{UserID: players[idx].ID, GameID: game.ID, Points: p.Score}
			err = pdb.Create(playerGame)
			if err != nil {
				return &MafiaBotPolemicaParserSaveMinimalGameStatisticPlayerGameError{Detail: err.Error()}
			}
		}
		return nil
	})
	return err
}
