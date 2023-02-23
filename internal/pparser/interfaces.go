package pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"io"
	"net/http"
)

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type PolemicaRequestInterface interface {
	Request(method, url string, body io.Reader, queryParams []*QueryParams) (*PolemicaResponse, error)
	GetCurrentUserID() int
}

type MafiaBotServiceInterface interface {
	SaveMinimalGameStatistic(MinimalGameStatistic) error
	GetLastGame() (*db.Game, error)
}

type PolemicaParserInterface interface {
	ParseGamesHistory(userID int, opts ...ParseGameHistoryOptionsParser) error
}
