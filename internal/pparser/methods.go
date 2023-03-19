package pparser

import (
	"bytes"
	"encoding/json"
	"github.com/GeorgeS1995/mafia_bot/internal"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/pparser"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type PolemicaApiClient struct {
	Requester           PolemicaRequestInterface
	mu                  sync.Mutex
	DBhandler           MafiaBotServiceInterface
	GameStatisticsRegex *regexp.Regexp
}

var GameStaticsRegex = ":game-data='(.+)'"

func NewPolemicaApiClient(cfg *pparser.MafiaBotPparserConfig, dbHandler MafiaBotServiceInterface) *PolemicaApiClient {
	match, _ := regexp.Compile(GameStaticsRegex)
	return &PolemicaApiClient{
		Requester:           NewPolemicaRequester(cfg),
		DBhandler:           dbHandler,
		GameStatisticsRegex: match,
	}
}

type PolemicaLoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *PolemicaApiClient) Login(username string, password string) error {
	loginPath := "mf/guest/login"

	b, err := json.Marshal(&PolemicaLoginBody{Username: username, Password: password})
	if err != nil {
		return &MafiaBotPolemicaParserLoginUnmarshalError{
			Detail: err.Error(),
		}
	}

	_, err = p.Requester.Request("POST", loginPath, bytes.NewBuffer(b), nil)
	if err != nil {
		return &MafiaBotPolemicaParserLoginResponselError{
			Detail: err.Error(),
		}
	}

	return nil
}

type MinimalPlayerGameStatistic struct {
	ID       string
	NickName string
	Score    float32
}

type MinimalGameStatistic struct {
	ID         string
	GameResult db.GameResult
	StartedAt  time.Time
	Players    [10]MinimalPlayerGameStatistic
}

type ParseGameHistoryOptions struct {
	Limit    int
	ToGameID string
}

type ParseGameHistoryOptionsParser func(o *ParseGameHistoryOptions)

func SetLimit(limit int) func(o *ParseGameHistoryOptions) {
	return func(o *ParseGameHistoryOptions) {
		o.Limit = limit
	}
}

func SetToGameID(gameID string) func(o *ParseGameHistoryOptions) {
	return func(o *ParseGameHistoryOptions) {
		o.ToGameID = gameID
	}
}

// Little helpers to determine the first error in goroutine

type GoroutineGameParseError struct {
	PageIdx        int
	GameParseError error
}

type GoroutineGameParseErrorArray []GoroutineGameParseError

func (a GoroutineGameParseErrorArray) GetFirstError() GoroutineGameParseError {
	minIdx := int(^uint(0) >> 1) // max integer
	var minGameParserError GoroutineGameParseError
	for _, e := range a {
		if e.PageIdx < minIdx {
			minGameParserError = e
			minIdx = e.PageIdx
		}
	}
	return minGameParserError
}

func (p *PolemicaApiClient) ParseGamesHistory(userID int, opts ...ParseGameHistoryOptionsParser) error {
	// kwargs patern for golang https://levelup.gitconnected.com/optional-function-parameter-pattern-in-golang-c1acc829307b
	options := &ParseGameHistoryOptions{
		30,
		"",
	}
	for _, opt := range opts {
		opt(options)
	}

	offset := 0

	limitOffsetQueryParams := []*QueryParams{
		{Param: "userId", Value: strconv.Itoa(userID)},
		{Param: "offset", Value: strconv.Itoa(offset)},
		{Param: "limit", Value: strconv.Itoa(options.Limit)},
	}
	gameHistoryPath := "cabinet/get"
	// When totalGameToParse == totalParsedGame than all goroutine done and we can leave the func
	totalGameToParse := 0
	totalParsedGame := 0
	pageParserCancelChan := make(chan bool)
	var errorsList GoroutineGameParseErrorArray
	defer func() {
		if totalGameToParse > 0 {
			<-pageParserCancelChan
		}
	}()
	for {
		polemicaGameHistoryResponse := &PolemicaGameHistoryResponse{}
		resp, err := p.Requester.Request("GET", gameHistoryPath, nil, limitOffsetQueryParams)
		if err != nil {
			return &MafiaBotPolemicaParserParseGamesHistoryResponseError{
				Detail:     err.Error(),
				QueryParam: limitOffsetQueryParams,
			}
		}

		err = json.Unmarshal(resp.Body, polemicaGameHistoryResponse)
		if err != nil {
			return &MafiaBotPolemicaParserParseGamesHistoryUnmarshalError{
				Detail: err.Error(),
			}
		}
		totalGameToParse += len(polemicaGameHistoryResponse.Rows)
		for idx, row := range polemicaGameHistoryResponse.Rows {
			gameId := row.Id
			if options.ToGameID == gameId {
				totalGameToParse -= len(polemicaGameHistoryResponse.Rows) - idx
				return nil
			}
			goroutineIdx := idx
			goroutineRow := row
			go func() {
				defer func() {
					totalParsedGame++
					// Unlock thread with main func if it's last game
					if totalParsedGame == totalGameToParse {
						pageParserCancelChan <- true
					}
				}()
				gameResult, goroutineErr := p.ParseGame(gameId)
				if goroutineErr != nil {
					errorsList = append(errorsList, GoroutineGameParseError{goroutineIdx, goroutineErr})
					return
				}
				// Add game started date
				date, goroutineErr := time.Parse(internal.PolemicaTimeFormat, goroutineRow.DateStart)
				if goroutineErr != nil {
					errorsList = append(errorsList, GoroutineGameParseError{goroutineIdx, goroutineErr})
					return
				}
				gameResult.StartedAt = date

				p.mu.Lock()
				goroutineErr = p.DBhandler.SaveMinimalGameStatistic(gameResult)
				if goroutineErr != nil {
					errorsList = append(errorsList, GoroutineGameParseError{goroutineIdx, goroutineErr})
				}
				p.mu.Unlock()
			}()
		}
		if len(errorsList) > 0 {
			p.mu.Lock()
			firstGoroutineGameParseError := errorsList.GetFirstError()
			p.mu.Unlock()
			return firstGoroutineGameParseError.GameParseError
		} else if len(polemicaGameHistoryResponse.Rows) < options.Limit {
			return nil
		}
		offset = offset + options.Limit
		limitOffsetQueryParams[1].Value = strconv.Itoa(offset)
	}
}

func (p *PolemicaApiClient) ParseGame(gameID string) (MinimalGameStatistic, error) {
	GameStatisticUrl := "game-statistics/" + gameID
	resp, err := p.Requester.Request("POST", GameStatisticUrl, nil, nil)
	if err != nil {
		return MinimalGameStatistic{}, &MafiaBotPolemicaParserParseGameResponseError{Detail: err.Error(), GameID: gameID}
	}
	gameStatisticsResponse := &GameStatisticsResponse{}
	stringBody := string(resp.Body)
	match := p.GameStatisticsRegex.FindStringSubmatch(stringBody)
	validatedBody := resp.Body
	if match != nil {
		validatedBody = []byte(match[1])
	}
	err = json.Unmarshal(validatedBody, gameStatisticsResponse)
	if err != nil {
		return MinimalGameStatistic{}, &MafiaBotPolemicaParserParseGameUnmarshalError{Detail: err.Error(), GameID: gameID}
	}
	minimalGameStatisticArray := [10]MinimalPlayerGameStatistic{}
	for idx, player := range gameStatisticsResponse.Players {
		minimalGameStatisticArray[idx] = MinimalPlayerGameStatistic{
			ID:       player.Id,
			NickName: player.Username,
			Score:    player.AchievementsSum.Points,
		}
	}
	gameResult, err := db.ToGameResult(gameStatisticsResponse.WinnerCode)
	if err != nil {
		return MinimalGameStatistic{}, &MafiaBotPolemicaParserParseGameEnumConvertationError{
			Detail: err.Error(),
		}
	}
	return MinimalGameStatistic{
		GameResult: gameResult,
		Players:    minimalGameStatisticArray,
		ID:         gameID,
	}, nil
}
