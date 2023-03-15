package test_pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	test_common "github.com/GeorgeS1995/mafia_bot/test/internal"
	"github.com/golang/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

func TestParseGameHistoryTaskOK(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	userID := rand.Intn(10000)
	ctrl := gomock.NewController(t)
	mockMafiaDB := test_common.NewMockMafiaDBInterface(ctrl)
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

	pparser.ParseGameHistoryTask(mockMafiaDB, mockParser, userID, test_common.TestTick, quitChan)

	time.Sleep(time.Duration(test_common.TestTick*2) * time.Millisecond)
	quitChan <- true
}

func TestParseGameHistoryTaskGetLastGameError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	userID := rand.Intn(10000)
	ctrl := gomock.NewController(t)
	mockMafiaDB := test_common.NewMockMafiaDBInterface(ctrl)
	mockParser := NewMockPolemicaParserInterface(ctrl)
	quitChan := make(chan bool)
	mockMafiaDB.EXPECT().GetLastGame().Return(&db.Game{}, &db.MafiaBotGetLastGameDriverError{}).Times(1)
	mockParser.EXPECT().ParseGamesHistory(gomock.Any(), gomock.Any()).Times(0)

	pparser.ParseGameHistoryTask(mockMafiaDB, mockParser, userID, test_common.TestTick, quitChan)

	time.Sleep(time.Duration(test_common.TestTick) * time.Millisecond)
	quitChan <- true
}
