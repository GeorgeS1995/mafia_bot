package test_pparser

import (
	"encoding/json"
	"errors"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/golang/mock/gomock"
	"io"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// Test p.ParseGame method
func TestPolemicaParserParseGameOK(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	gameId := strconv.Itoa(rand.Intn(100001))
	gameStatisticResponse := RandomGameStatisticsResponse(gameId)
	b, _ := json.Marshal(gameStatisticResponse)
	ctrl := gomock.NewController(t)
	m := NewMockPolemicaRequestInterface(ctrl)
	m.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&pparser.PolemicaResponse{Body: b, StatusCode: 200}, nil)
	pParser := pparser.PolemicaApiClient{
		Requester: m,
	}

	resp, err := pParser.ParseGame(gameId)

	if err != nil {
		t.Fatalf("ParseGame unexpected error: %s", err.Error())
	}
	expectedGameResult, err := db.ToGameResult(gameStatisticResponse.WinnerCode)
	if err != nil {
		t.Fatalf("Error while convert winner code to enum: %s", err.Error())
	}
	if resp.GameResult != expectedGameResult {
		t.Fatalf("ParseGame winner codes aren't equal. \nExpected: %s\nAsserted %s", expectedGameResult, resp.GameResult)
	}
	for idx, p := range gameStatisticResponse.Players {
		if p.Id != resp.Players[idx].ID {
			t.Fatalf("ParseGame player with index %d incorrect unmarshalled. \nExpected: %+v\nAsserted %+v", idx, p, resp.Players[idx])
		}
	}
}

func TestPolemicaParserParseGameResponseError(t *testing.T) {
	requestError := errors.New("")
	rand.Seed(time.Now().UnixNano())
	gameId := strconv.Itoa(rand.Intn(100001))
	ctrl := gomock.NewController(t)
	m := NewMockPolemicaRequestInterface(ctrl)
	m.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&pparser.PolemicaResponse{}, requestError)
	pParser := pparser.PolemicaApiClient{
		Requester: m,
	}

	_, err := pParser.ParseGame(gameId)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserParseGameResponseError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestPolemicaParserParseGameUnmarshalError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	gameId := strconv.Itoa(rand.Intn(100001))
	ctrl := gomock.NewController(t)
	m := NewMockPolemicaRequestInterface(ctrl)
	m.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&pparser.PolemicaResponse{Body: nil, StatusCode: 200}, nil)
	pParser := pparser.PolemicaApiClient{
		Requester: m,
	}

	_, err := pParser.ParseGame(gameId)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserParseGameUnmarshalError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestPolemicaParserParseGameEnumConvertationError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	gameId := strconv.Itoa(rand.Intn(100001))
	gameStatisticResponse := RandomGameStatisticsResponse(gameId)
	gameStatisticResponse.WinnerCode = 3
	b, _ := json.Marshal(gameStatisticResponse)
	ctrl := gomock.NewController(t)
	m := NewMockPolemicaRequestInterface(ctrl)
	m.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&pparser.PolemicaResponse{Body: b, StatusCode: 200}, nil)
	pParser := pparser.PolemicaApiClient{
		Requester: m,
	}

	_, err := pParser.ParseGame(gameId)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserParseGameEnumConvertationError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

// Test p.ParseGamesHistory method
func TestPolemicaParserParseGamesHistoryOK(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	polemicaUserId := rand.Intn(100001)
	ctrl := gomock.NewController(t)
	requestMock := NewMockPolemicaRequestInterface(ctrl)
	// api response mock
	isSecondPage := false
	requestMock.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(method, url string, body io.Reader, queryParams []*pparser.QueryParams) (*pparser.PolemicaResponse, error) {
			response := &pparser.PolemicaResponse{StatusCode: 200}
			var marshaled []byte
			if method == "GET" {
				rows := []pparser.PolemicaGameHistoryResponseRow{
					{Id: strconv.Itoa(rand.Intn(100001)), DateStart: pparser.PolemicaTimeFormat},
				}
				if !isSecondPage {
					rows = append(rows, pparser.PolemicaGameHistoryResponseRow{Id: strconv.Itoa(rand.Intn(100001)), DateStart: pparser.PolemicaTimeFormat})
					isSecondPage = true
				}
				toMarshal := &pparser.PolemicaGameHistoryResponse{
					Rows: rows,
				}
				marshaled, _ = json.Marshal(toMarshal)

			} else if method == "POST" {
				toMarshal := RandomGameStatisticsResponse(strconv.Itoa(rand.Intn(100001)))
				marshaled, _ = json.Marshal(toMarshal)

			}
			response.Body = marshaled
			return response, nil
		}).Times(5)
	// Saver mock
	dbMock := NewMockMafiaBotServiceInterface(ctrl)
	dbMock.EXPECT().SaveMinimalGameStatistic(gomock.Any()).Times(3)
	// Create test PolemicaApiClient
	pParser := pparser.PolemicaApiClient{
		Requester: requestMock,
		DBhandler: dbMock,
	}

	err := pParser.ParseGamesHistory(polemicaUserId, pparser.SetLimit(2))

	if err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
	}
}

func TestPolemicaParserParseGamesHistoryOKStopByGameId(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	polemicaUserId := rand.Intn(100001)
	ctrl := gomock.NewController(t)
	requestMock := NewMockPolemicaRequestInterface(ctrl)
	// api response mock
	stopGameID := strconv.Itoa(rand.Intn(100001))
	isSecondPage := false
	requestMock.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(method, url string, body io.Reader, queryParams []*pparser.QueryParams) (*pparser.PolemicaResponse, error) {
			response := &pparser.PolemicaResponse{StatusCode: 200}
			var marshaled []byte
			if method == "GET" {
				gameId := strconv.Itoa(rand.Intn(100001))
				if isSecondPage {
					gameId = stopGameID
				}
				toMarshal := &pparser.PolemicaGameHistoryResponse{
					Rows: []pparser.PolemicaGameHistoryResponseRow{
						{Id: strconv.Itoa(rand.Intn(100001)), DateStart: pparser.PolemicaTimeFormat},
						{Id: gameId, DateStart: pparser.PolemicaTimeFormat},
					},
				}
				marshaled, _ = json.Marshal(toMarshal)
				isSecondPage = true
			} else if method == "POST" {
				toMarshal := RandomGameStatisticsResponse(strconv.Itoa(rand.Intn(100001)))
				marshaled, _ = json.Marshal(toMarshal)
			}
			response.Body = marshaled
			return response, nil
		}).Times(5)
	// Saver mock
	dbMock := NewMockMafiaBotServiceInterface(ctrl)
	dbMock.EXPECT().SaveMinimalGameStatistic(gomock.Any()).Times(3)
	// Create test PolemicaApiClient
	pParser := pparser.PolemicaApiClient{
		Requester: requestMock,
		DBhandler: dbMock,
	}

	err := pParser.ParseGamesHistory(polemicaUserId, pparser.SetLimit(2), pparser.SetToGameID(stopGameID))

	if err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
	}
}

func TestPolemicaParserParseGamesHistoryResponseError(t *testing.T) {
	polemicaUserId := rand.Intn(100001)
	ctrl := gomock.NewController(t)
	requestMock := NewMockPolemicaRequestInterface(ctrl)
	requestMock.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&pparser.PolemicaResponse{}, errors.New(""))
	pParser := pparser.PolemicaApiClient{
		Requester: requestMock,
	}

	err := pParser.ParseGamesHistory(polemicaUserId, pparser.SetLimit(2))

	if _, ok := err.(*pparser.MafiaBotPolemicaParserParseGamesHistoryResponseError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestPolemicaParserParseGamesHistoryUnmarshalError(t *testing.T) {
	polemicaUserId := rand.Intn(100001)
	ctrl := gomock.NewController(t)
	requestMock := NewMockPolemicaRequestInterface(ctrl)
	requestMock.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&pparser.PolemicaResponse{Body: nil, StatusCode: 200}, nil)
	pParser := pparser.PolemicaApiClient{
		Requester: requestMock,
	}

	err := pParser.ParseGamesHistory(polemicaUserId, pparser.SetLimit(2))

	if _, ok := err.(*pparser.MafiaBotPolemicaParserParseGamesHistoryUnmarshalError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestPolemicaParserParseGamesHistoryGoroutineError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	polemicaUserId := rand.Intn(100001)
	ctrl := gomock.NewController(t)
	requestMock := NewMockPolemicaRequestInterface(ctrl)
	dbMock := NewMockMafiaBotServiceInterface(ctrl)
	dbMock.EXPECT().SaveMinimalGameStatistic(gomock.Any()).Return(nil).AnyTimes()
	// api response mock
	isFirstGame := false
	gameId := strconv.Itoa(rand.Intn(100001))
	requestMock.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(method, url string, body io.Reader, queryParams []*pparser.QueryParams) (*pparser.PolemicaResponse, error) {
			response := &pparser.PolemicaResponse{}
			var mockError error
			if method == "GET" {
				to_marshal := &pparser.PolemicaGameHistoryResponse{
					Rows: []pparser.PolemicaGameHistoryResponseRow{
						{Id: gameId},
						{Id: strconv.Itoa(rand.Intn(100001))},
					},
				}
				marshaled, _ := json.Marshal(to_marshal)
				response.StatusCode = 200
				response.Body = marshaled
			} else if method == "POST" {
				pasedGameId := strconv.Itoa(rand.Intn(100001))
				if isFirstGame {
					pasedGameId = gameId
					isFirstGame = false
				}
				_ = RandomGameStatisticsResponse(pasedGameId)
				response.StatusCode = 400
				mockError = &pparser.MafiaBotPolemicaParserNewRequestError{}
			}
			return response, mockError
		}).AnyTimes()

	pParser := pparser.PolemicaApiClient{
		Requester: requestMock,
		DBhandler: dbMock,
	}
	err := pParser.ParseGamesHistory(polemicaUserId, pparser.SetLimit(2))

	if typedError, ok := err.(*pparser.MafiaBotPolemicaParserParseGameResponseError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	} else if typedError.GameID != gameId {
		t.Fatalf("Error from wrong gogoutine")
	}
}

// goroutineGameParseErrorArray tests
func TestGetFirstError(t *testing.T) {
	errorList := pparser.GoroutineGameParseErrorArray{
		{3, errors.New("")},
		{1, errors.New("")},
		{2, errors.New("")},
	}

	firstIdxError := errorList.GetFirstError()

	if firstIdxError.PageIdx != 1 {
		t.Fatalf("Wrong error idx: %d", firstIdxError.PageIdx)
	}
}
