package pparser

import (
	"io"
	"net/http"
)

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type PolemicaRequestInterface interface {
	Request(method, url string, body io.Reader, queryParams []*QueryParams) (*PolemicaResponse, error)
}

type MafiaBotDBInterface interface {
	SaveMinimalGameStatistic(MinimalGameStatistic) error
}
