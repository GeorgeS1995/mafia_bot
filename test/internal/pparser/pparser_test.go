package test_pparser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/golang/mock/gomock"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestSetAuthCookie(t *testing.T) {
	// 1 unchanged 1 changed 1 new cookie
	WillNotModifiedValue := test.RandStringRunes(3)
	WillNotModidiedName := "WillNotModified"
	WillChangedName := "WillChanged"
	originAuthCookie := []*http.Cookie{{
		Name:  WillNotModidiedName,
		Value: WillNotModifiedValue,
	},
		{
			Name:  WillChangedName,
			Value: test.RandStringRunes(3),
		},
	}
	newWillChangedValue := test.RandStringRunes(4)
	NewCookieValue := test.RandStringRunes(4)
	NewCookieName := "NewCookie"
	newRawCookies := []string{
		fmt.Sprintf("%s=%s; expires=Sun, 08-Jan-2023 14:02:39 GMT; Max-Age=86400; path=/; HttpOnly", WillChangedName, newWillChangedValue),
		fmt.Sprintf("%s=%s; expires=Sun, 08-Jan-2023 14:02:39 GMT; Max-Age=86400; path=/; HttpOnly", NewCookieName, NewCookieValue),
	}
	header := http.Header{"Set-Cookie": newRawCookies}
	pParser := pparser.PolemicaRequester{
		AuthCookie: originAuthCookie,
	}

	pParser.SetAuthCookie(header)

	for _, cookie := range pParser.AuthCookie {
		if cookie.Name == WillNotModidiedName && cookie.Value != WillNotModifiedValue {
			t.Fatalf("%s cookie was modified", WillNotModidiedName)
		} else if cookie.Name == WillChangedName && cookie.Value != newWillChangedValue {
			t.Fatalf("%s cookie was not modified", WillChangedName)
		}
	}
	newCookie := pParser.AuthCookie[len(pParser.AuthCookie)-1]
	if newCookie.Name != NewCookieName || newCookie.Value != NewCookieValue {
		t.Fatalf("%s cookie was not added to AuthCookie", NewCookieName)
	}
}

func TestPolemicaRequestNewRequestError(t *testing.T) {
	pParser := pparser.PolemicaRequester{}
	method := "GET>"
	url := ""
	ctrl := gomock.NewController(t)
	m := NewMockHttpClientInterface(ctrl)

	_, err := pParser.PolemicaRequest(method, url, nil, m, nil)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserNewRequestError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestPolemicaClientDoError(t *testing.T) {
	pParser := pparser.PolemicaRequester{}
	method := "GET"
	url := ""
	ctrl := gomock.NewController(t)
	m := NewMockHttpClientInterface(ctrl)
	m.EXPECT().Do(gomock.Any()).Return(&http.Response{}, errors.New(""))

	_, err := pParser.PolemicaRequest(method, url, nil, m, nil)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserRequestConnectionError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

// Always read with error for testing purpose
type ReadCloserError struct{}

func (r *ReadCloserError) Read(p []byte) (n int, err error) {
	return 0, errors.New("")
}
func (r *ReadCloserError) Close() error {
	return nil
}

func TestPolemicaParserResponseBodyParsingError(t *testing.T) {
	pParser := pparser.PolemicaRequester{}
	method := "GET"
	url := ""
	ctrl := gomock.NewController(t)
	m := NewMockHttpClientInterface(ctrl)
	m.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: &ReadCloserError{}}, nil)

	_, err := pParser.PolemicaRequest(method, url, nil, m, nil)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserResponseBodyParsingError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestPolemicaParserServerResponseError(t *testing.T) {
	pParser := pparser.PolemicaRequester{}
	method := "GET"
	url := ""
	rand.Seed(time.Now().UnixNano())
	ctrl := gomock.NewController(t)
	m := NewMockHttpClientInterface(ctrl)
	m.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte{255})), StatusCode: rand.Intn(599-300+1) + 300}, nil)

	_, err := pParser.PolemicaRequest(method, url, nil, m, nil)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserServerResponseError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

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
			marshaled := []byte{}
			if method == "GET" {
				rows := []pparser.PolemicaGameHistoryResponseRow{
					{Id: strconv.Itoa(rand.Intn(100001)), DateStart: pparser.PolemicaTimeFormat},
				}
				if !isSecondPage {
					rows = append(rows, pparser.PolemicaGameHistoryResponseRow{Id: strconv.Itoa(rand.Intn(100001)), DateStart: pparser.PolemicaTimeFormat})
					isSecondPage = true
				}
				to_marshal := &pparser.PolemicaGameHistoryResponse{
					Rows: rows,
				}
				marshaled, _ = json.Marshal(to_marshal)

			} else if method == "POST" {
				to_marshal := RandomGameStatisticsResponse(strconv.Itoa(rand.Intn(100001)))
				marshaled, _ = json.Marshal(to_marshal)

			}
			response.Body = marshaled
			return response, nil
		}).Times(5)
	// Saver mock
	dbMock := NewMockMafiaBotDBInterface(ctrl)
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
			marshaled := []byte{}
			if method == "GET" {
				gameId := strconv.Itoa(rand.Intn(100001))
				if isSecondPage {
					gameId = stopGameID
				}
				to_marshal := &pparser.PolemicaGameHistoryResponse{
					Rows: []pparser.PolemicaGameHistoryResponseRow{
						{Id: strconv.Itoa(rand.Intn(100001)), DateStart: pparser.PolemicaTimeFormat},
						{Id: gameId, DateStart: pparser.PolemicaTimeFormat},
					},
				}
				marshaled, _ = json.Marshal(to_marshal)
				isSecondPage = true
			} else if method == "POST" {
				to_marshal := RandomGameStatisticsResponse(strconv.Itoa(rand.Intn(100001)))
				marshaled, _ = json.Marshal(to_marshal)
			}
			response.Body = marshaled
			return response, nil
		}).Times(5)
	// Saver mock
	dbMock := NewMockMafiaBotDBInterface(ctrl)
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
