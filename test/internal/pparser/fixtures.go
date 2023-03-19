package test_pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func RandomGameStatisticsResponse(gameID string) pparser.GameStatisticsResponse {
	rand.Seed(time.Now().UnixNano())
	var players []pparser.GameStatisticsPlayerResponse
	for i := 0; i < 10; i++ {
		players = append(players, pparser.GameStatisticsPlayerResponse{
			Id:       strconv.Itoa(rand.Intn(100001)),
			Username: test.RandStringRunes(5),
			AchievementsSum: pparser.GameStatisticsAchievementsSum{
				Points: rand.Float32(),
			},
		})
	}
	return pparser.GameStatisticsResponse{
		Players:    players,
		WinnerCode: rand.Intn(2) + 1,
		Id:         gameID,
	}
}

func GetTestPolemicaApiClient(r pparser.PolemicaRequestInterface, db pparser.MafiaBotServiceInterface) pparser.PolemicaApiClient {
	match, _ := regexp.Compile(pparser.GameStaticsRegex)
	return pparser.PolemicaApiClient{
		Requester:           r,
		DBhandler:           db,
		GameStatisticsRegex: match,
	}
}
