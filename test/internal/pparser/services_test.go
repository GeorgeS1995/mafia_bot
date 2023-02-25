package test_pparser

import (
	"database/sql"
	"errors"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	mockdb "github.com/GeorgeS1995/mafia_bot/test/internal"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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
	pdb := &pparser.MafiaParserServiceHandler{MafiaDBInterface: &db.MafiaDB{suite.Tx}}
	players := [10]pparser.MinimalPlayerGameStatistic{}
	existingPlayer := &db.User{}
	for idx := range players {
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

func TestMafiaParserServiceHandlerUpdateOrCreateError(t *testing.T) {
	minimalGameStatistic := pparser.MinimalGameStatistic{}
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockMafiaDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(
		func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
			return fc(&gorm.DB{})
		})
	mockMafiaDB.EXPECT().UpdateOrCreateUser(gomock.Any()).Return(nil, errors.New(""))
	mockMafiaDB.EXPECT().Create(gomock.Any()).Times(0)
	pdb := &pparser.MafiaParserServiceHandler{MafiaDBInterface: mockMafiaDB}

	err := pdb.SaveMinimalGameStatistic(minimalGameStatistic)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserSaveMinimalGameStatisticUserError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestMafiaParserServiceHandlerGameCreateError(t *testing.T) {
	minimalGameStatistic := pparser.MinimalGameStatistic{}
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockMafiaDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(
		func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
			return fc(&gorm.DB{})
		})
	mockMafiaDB.EXPECT().UpdateOrCreateUser(gomock.Any()).Return(&db.User{}, nil).Times(10)
	mockMafiaDB.EXPECT().Create(gomock.Any()).Return(errors.New(""))
	pdb := &pparser.MafiaParserServiceHandler{MafiaDBInterface: mockMafiaDB}

	err := pdb.SaveMinimalGameStatistic(minimalGameStatistic)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserSaveMinimalGameStatisticGameError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestMafiaParserServiceHandlerPlayerGameError(t *testing.T) {
	minimalGameStatistic := pparser.MinimalGameStatistic{}
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockMafiaDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(
		func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
			return fc(&gorm.DB{})
		})
	mockMafiaDB.EXPECT().UpdateOrCreateUser(gomock.Any()).Return(&db.User{}, nil).Times(10)
	createCallCount := 0
	mockMafiaDB.EXPECT().Create(gomock.Any()).DoAndReturn(
		func(value interface{}) error {
			var err error
			if createCallCount == 0 {
				createCallCount++
			} else {
				err = errors.New("")
			}
			return err
		},
	).Times(2)
	pdb := &pparser.MafiaParserServiceHandler{MafiaDBInterface: mockMafiaDB}

	err := pdb.SaveMinimalGameStatistic(minimalGameStatistic)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserSaveMinimalGameStatisticPlayerGameError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}
