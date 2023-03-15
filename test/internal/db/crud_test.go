package test_db

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type DBPackageTestSuite struct {
	test.DBTestSuite
}

func TestDBPackage(t *testing.T) {
	suite.Run(t, new(DBPackageTestSuite))
}

func (suite *DBPackageTestSuite) TestUpdateOrCreateUserExistingUserOK() {
	defer suite.Tx.Rollback()
	polemicaID := test.RandStringRunes(3)
	polemicaNickName := test.RandStringRunes(3)
	polemicaNickNameNew := test.RandStringRunes(3)
	existingUser := &db.User{
		PolemicaId:       polemicaID,
		PolemicaNickName: polemicaNickName,
	}
	suite.Tx.Create(existingUser)
	creatAt := existingUser.CreatedAt
	mafiaDB := &db.MafiaDB{Db: suite.Tx}

	user, err := mafiaDB.UpdateOrCreateUser(&db.User{
		PolemicaId:       polemicaID,
		PolemicaNickName: polemicaNickNameNew,
	})

	timeNow := time.Now()
	if err != nil {
		suite.T().Fatal("Update user error, ", err)
	}
	assertedUser := &db.User{}
	var count int64
	suite.Tx.Model(&db.User{}).Where("polemica_id = ?", polemicaID).Count(&count).First(assertedUser)
	assert.Equal(suite.T(), int(count), 1)
	assert.Equal(suite.T(), user.PolemicaId, assertedUser.PolemicaId)
	assert.Equal(suite.T(), assertedUser.PolemicaNickName, polemicaNickNameNew)
	assert.Equal(suite.T(), assertedUser.CreatedAt.Round(time.Second), creatAt.Round(time.Second))
	assert.Equal(suite.T(), assertedUser.UpdatedAt.Round(time.Minute), timeNow.Round(time.Minute))
}

func (suite *DBPackageTestSuite) TestUpdateOrCreateUserCreateOK() {
	defer suite.Tx.Rollback()
	polemicaID := test.RandStringRunes(3)
	polemicaNickName := test.RandStringRunes(3)
	mafiaDB := &db.MafiaDB{Db: suite.Tx}

	user, err := mafiaDB.UpdateOrCreateUser(&db.User{
		PolemicaId:       polemicaID,
		PolemicaNickName: polemicaNickName,
	})

	timeNow := time.Now()
	if err != nil {
		suite.T().Fatal("Create user error, ", err)
	}
	assertedUser := &db.User{}
	var count int64
	suite.Tx.Model(&db.User{}).Where("polemica_id = ?", polemicaID).Count(&count).First(assertedUser)
	assert.Equal(suite.T(), int(count), 1)
	assert.Equal(suite.T(), user.PolemicaId, assertedUser.PolemicaId)
	assert.Equal(suite.T(), assertedUser.PolemicaNickName, polemicaNickName)
	assert.Equal(suite.T(), assertedUser.UpdatedAt.Round(time.Minute), timeNow.Round(time.Minute))
	assert.Equal(suite.T(), assertedUser.CreatedAt.Round(time.Minute), timeNow.Round(time.Minute))
}

func (suite *DBPackageTestSuite) TestGetLastGame() {
	defer suite.Tx.Rollback()
	suite.Tx.Create(&db.Game{StartedAt: time.Now(), Winner: db.Draw})
	lastGameTime := time.Now().Add(time.Hour)
	expectedlastGame := &db.Game{StartedAt: lastGameTime, Winner: db.Draw}
	suite.Tx.Create(expectedlastGame)
	mafiaDB := &db.MafiaDB{Db: suite.Tx}

	lastGame, err := mafiaDB.GetLastGame()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedlastGame.ID, lastGame.ID)
	assert.Equal(suite.T(), lastGame.StartedAt.Round(time.Second), lastGameTime.Round(time.Second))
}

func (suite *DBPackageTestSuite) TestGetLastGameEmptyDB() {
	defer suite.Tx.Rollback()
	mafiaDB := &db.MafiaDB{Db: suite.Tx}

	_, err := mafiaDB.GetLastGame()
	if _, ok := err.(*db.MafiaBotGetLastGameEmptyDBError); !ok {
		suite.T().Fatalf("Wrong error type: %s", err)
	}
}

func (suite *DBPackageTestSuite) TestGetDailyStatisticOK() {
	defer suite.Tx.Rollback()
	for _, switchHour := range [3]int{0, 1, 23} {
		suite.SetupTest()
		suite.Run(fmt.Sprintf("Current switch hour %d", switchHour), func() {
			defer func() {
				suite.Tx.Rollback()
			}()
			mafiaDB := &db.MafiaDB{Db: suite.Tx}
			timeNow := time.Now()
			_, offset := timeNow.Zone()
			offsetHour := offset / 3600
			// Already send game
			CreateRandomGame(mafiaDB, db.Game{IsSendedToDiscord: true})
			// Game after time switch
			CreateRandomGame(mafiaDB, db.Game{StartedAt: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), switchHour+1, 0, 0, 0, timeNow.Location())})
			// Today games
			firstDayMVP := db.MVP{
				NickName: test.RandStringRunes(5),
			}
			todayGameIDs := []uuid.UUID{}
			for i := 0; i < 2; i++ {
				firstDayMVP.Score = 0.5 + float32(i)
				duration := switchHour - 1 - i - offsetHour
				if duration < 0 {
					duration = duration * -1
				}
				gameStartedAt := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), switchHour, 0, 0, 0, timeNow.Location()).Add(-time.Hour * time.Duration(duration))
				game := CreateRandomGame(mafiaDB, db.Game{StartedAt: gameStartedAt, Winner: db.MafiaWin}, firstDayMVP)
				todayGameIDs = append(todayGameIDs, game.ID)
			}
			// Yesterday games
			secondDayMVP := db.MVP{
				NickName: test.RandStringRunes(5),
				Score:    2,
			}
			firstDayMVP.Score = 1
			secondDayGame := CreateRandomGame(mafiaDB, db.Game{StartedAt: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()-1, switchHour, 0, 0, 0, timeNow.Location()), Winner: db.CityWin}, firstDayMVP, secondDayMVP)

			lastGameList, err := mafiaDB.GetDailyStatistic(switchHour)

			if err != nil {
				suite.T().Fatalf("Unexpected error: %s", err)
			}
			assert.Equal(suite.T(), len(lastGameList), 2)
			utcLoc, _ := time.LoadLocation("UTC")
			utcNow := timeNow.In(utcLoc)
			// Additional day offset for edge cases
			var dayOffset int
			if switchHour <= 1 {
				dayOffset = -1
			}
			for idx, gameStatic := range lastGameList {
				if idx == 0 {
					assert.Equal(suite.T(), gameStatic.Date, time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day()-1+dayOffset, 0, 0, 0, 0, utcNow.Location()))
					assert.Equal(suite.T(), gameStatic.GameCount, 1)
					assert.Equal(suite.T(), gameStatic.MafiaWins, 0)
					assert.Equal(suite.T(), gameStatic.CityWins, 1)
					assert.Equal(suite.T(), gameStatic.MVP.Score, float32(2))
					assert.Equal(suite.T(), gameStatic.MVP.NickName, secondDayMVP.NickName)
					assert.Equal(suite.T(), gameStatic.PlayedGamesID[0], secondDayGame.ID)
				} else {
					assert.Equal(suite.T(), gameStatic.Date, time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day()+dayOffset, 0, 0, 0, 0, utcNow.Location()))
					assert.Equal(suite.T(), gameStatic.GameCount, 2)
					assert.Equal(suite.T(), gameStatic.MafiaWins, 2)
					assert.Equal(suite.T(), gameStatic.CityWins, 0)
					assert.Equal(suite.T(), gameStatic.MVP.Score, float32(1.5))
					assert.Equal(suite.T(), gameStatic.MVP.NickName, firstDayMVP.NickName)
					for _, gameID := range todayGameIDs {
						assert.Contains(suite.T(), gameStatic.PlayedGamesID, gameID)
					}
				}
			}
		})
	}
}

func (suite *DBPackageTestSuite) TestMarkGamesAsSentOK() {
	defer suite.Tx.Rollback()
	mafiaDB := &db.MafiaDB{Db: suite.Tx}
	gameIDS := []uuid.UUID{}
	for i := 0; i < 2; i++ {
		gameToCreate := &db.Game{IsSendedToDiscord: false, Winner: db.Draw}
		suite.Tx.Create(gameToCreate)
		gameIDS = append(gameIDS, gameToCreate.ID)
	}

	err := mafiaDB.MarkGamesAsSent(gameIDS)

	if err != nil {
		log.Fatalf("Unexpected error in TestMarkGamesAsSentOK: %s", err.Error())
	}
	games := []*db.Game{}
	suite.Tx.Find(&games, gameIDS)
	for _, game := range games {
		assert.Contains(suite.T(), gameIDS, game.ID)
		assert.Equal(suite.T(), true, game.IsSendedToDiscord)
	}
}
