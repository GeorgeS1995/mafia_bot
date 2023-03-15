package test_db

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/test"
	"math"
	"math/rand"
	"time"
)

// TODO Should we handle error in fixture?
func CreateRandomGame(mdb *db.MafiaDB, game db.Game, opts ...db.MVP) db.Game {
	rand.Seed(time.Now().UnixNano())
	if game.PolemicaId == "" {
		game.PolemicaId = test.RandStringRunes(3)
	}
	defaultTime := time.Time{}
	if game.StartedAt == defaultTime {
		game.StartedAt = time.Now()
	}
	if game.Winner == "" {
		game.Winner = db.Draw
	}
	mdb.Create(&game)
	var maxScore float32
	for _, u := range opts {
		newUser := &db.User{
			PolemicaNickName: u.NickName,
			PolemicaId:       test.RandStringRunes(3),
		}
		mdb.Create(newUser)
		mdb.Create(&db.PlayerGame{
			UserID: newUser.ID,
			GameID: game.ID,
			Points: u.Score,
		})
		if maxScore < u.Score {
			maxScore = u.Score
		}
	}
	for i := 1; i <= 10-len(opts); i++ {
		newUser := &db.User{
			PolemicaNickName: test.RandStringRunes(3),
			PolemicaId:       test.RandStringRunes(3),
		}
		mdb.Create(newUser)
		maxPoints := 10
		if maxScore > 0 {
			maxPoints = int(math.Round(float64(maxScore)*10.0 - 1))
		}
		points := float32(rand.Intn(maxPoints) / 10)
		mdb.Create(&db.PlayerGame{
			UserID: newUser.ID,
			GameID: game.ID,
			Points: points,
		})
	}
	return game
}
