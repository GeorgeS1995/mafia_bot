package test_db

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"math/rand"
	"testing"
	"time"
)

type GameResultMap struct {
	WinnerCode       int
	GameResultString db.GameResult
}

func TestToGameResultOK(t *testing.T) {
	winnerCodes := []GameResultMap{
		{0, db.Draw},
		{1, db.CityWin},
		{2, db.MafiaWin},
	}
	for _, wc := range winnerCodes {
		t.Run(fmt.Sprintf("Test convertion code %d to GameResult", wc.WinnerCode), func(t *testing.T) {

			gr, err := db.ToGameResult(wc.WinnerCode)

			if err != nil {
				t.Fatal(fmt.Sprintf("Unexpect error: %s", err.Error()))
			}
			if gr != wc.GameResultString {
				t.Fatal(fmt.Sprintf("Unexpected result. Expected: %s, Got: %s", wc.GameResultString, gr))
			}
		})
	}
}

func TestToGameResultError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	wrongCode := rand.Intn(10) + 3

	_, err := db.ToGameResult(wrongCode)

	if _, ok := err.(*db.MafiaBotToGameResultError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}
