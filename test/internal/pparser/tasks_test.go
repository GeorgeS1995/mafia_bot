package test_pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	mockdb "github.com/GeorgeS1995/mafia_bot/test/internal"
	"github.com/golang/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

const testTick = 10

func TestParseGameHistoryTaskOK(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	userID := rand.Intn(10000)
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockParser := NewMockPolemicaParserInterface(ctrl)
	quitChan := make(chan bool)
	dbCallCount := 0
	mockMafiaDB.EXPECT().GetLastGame().DoAndReturn(func() (*db.Game, error) {
		if dbCallCount == 0 {
			dbCallCount++
			return nil, &db.MafiaBotGetLastGameEmptyDBError{}
		} else {
			return &db.Game{PolemicaId: test.RandStringRunes(3)}, nil
		}
	}).Times(2)
	mockParser.EXPECT().ParseGamesHistory(gomock.Any(), gomock.Any()).Times(2)

	pparser.ParseGameHistoryTask(mockMafiaDB, mockParser, userID, testTick, quitChan)

	time.Sleep(time.Duration(testTick*2) * time.Millisecond)
	quitChan <- true
}

func TestParseGameHistoryTaskGetLastGameError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	userID := rand.Intn(10000)
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockParser := NewMockPolemicaParserInterface(ctrl)
	quitChan := make(chan bool)
	mockMafiaDB.EXPECT().GetLastGame().Return(&db.Game{}, &db.MafiaBotGetLastGameDriverError{}).Times(1)
	mockParser.EXPECT().ParseGamesHistory(gomock.Any(), gomock.Any()).Times(0)

	pparser.ParseGameHistoryTask(mockMafiaDB, mockParser, userID, testTick, quitChan)

	time.Sleep(time.Duration(testTick) * time.Millisecond)
	quitChan <- true
}
