package test_pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type PParserTestSuite struct {
	test.DBTestSuite
}

func TestDBPackage(t *testing.T) {
	suite.Run(t, new(PParserTestSuite))
}

func (suite *PParserTestSuite) TestSaveMinimalGameStatisticOK() {
	defer suite.Tx.Rollback()
	pdb := &pparser.ParserMafiaDB{db.MafiaDB{suite.Tx}}
	players := [10]pparser.MinimalPlayerGameStatistic{}
	existingPlayer := &db.User{}
	for idx, _ := range players {
		polemicaId := test.RandStringRunes(5)
		polemicaNickName := test.RandStringRunes(3)
		if idx == 0 {
			existingPlayer.PolemicaId = polemicaId
			existingPlayer.PolemicaNickName = polemicaNickName
			suite.Tx.Create(existingPlayer)
		}
		players[idx].ID = polemicaId
		players[idx].NickName = polemicaNickName
		players[idx].Score = float32(idx) * 0.1
	}
	minimalGameStatistic := pparser.MinimalGameStatistic{
		Players:    players,
		ID:         test.RandStringRunes(3),
		GameResult: db.CityWin,
		StartedAt:  time.Now().Add(-time.Hour * 24),
	}

	err := pdb.SaveMinimalGameStatistic(minimalGameStatistic)

	if err != nil {
		suite.T().Fatal("Unexpected error while saving game statistic, ", err)
	}
	var userCount int64
	suite.Tx.Model(&db.User{}).Count(&userCount)
	assert.Equal(suite.T(), int(userCount), 10)
	game := &db.Game{}
	suite.Tx.Where("polemica_id = ?", minimalGameStatistic.ID).First(&game)
	assert.Equal(suite.T(), game.StartedAt.Round(time.Minute), minimalGameStatistic.StartedAt.Round(time.Minute))
	assert.Equal(suite.T(), game.IsSendedToDiscord, false)
	assert.Equal(suite.T(), game.Winner, minimalGameStatistic.GameResult)
	for _, p := range players {
		player := &db.User{}
		suite.Tx.Where("polemica_id = ?", p.ID).First(&player)
		assert.Equal(suite.T(), player.PolemicaNickName, p.NickName)
		playerGame := &db.PlayerGame{}
		suite.Tx.Where(&db.PlayerGame{UserID: player.ID, GameID: game.ID}).First(&playerGame)
		assert.Equal(suite.T(), playerGame.Points, p.Score)
	}
}
