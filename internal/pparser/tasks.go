package pparser

import (
	dbPkg "github.com/GeorgeS1995/mafia_bot/internal/db"
	"log"
	"time"
)

func ParseGameHistoryTask(mdb dbPkg.MafiaDBInterface, parser PolemicaParserInterface, userID int, tickerDelay int, quit chan bool) {
	ticker := time.NewTicker(time.Duration(tickerDelay) * time.Millisecond)
	go func() {
		for {
		SB:
			select {
			case <-ticker.C:
				lastParsedGame, err := mdb.GetLastGame()
				switch err.(type) {
				case *dbPkg.MafiaBotGetLastGameDriverError:
					log.Printf("Error in GetLastGame func: %s", err.Error())
					break SB
				case *dbPkg.MafiaBotGetLastGameEmptyDBError:
					lastParsedGame = &dbPkg.Game{}
				}
				err = parser.ParseGamesHistory(userID, SetToGameID(lastParsedGame.PolemicaId))
				if err != nil {
					log.Printf("Error while parsing game history: %s", err.Error())
				}
			case <-quit:
				return
			}
		}
	}()
}
