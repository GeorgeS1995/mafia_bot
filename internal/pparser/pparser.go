package pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/pparser"
	"io"
	"net/http"
	"strings"
)

var client http.Client

type PolemicaRequester struct {
	Url               string
	Csrf              string
	CsrfInitialCookie string
	AuthCookie        []*http.Cookie
}

// Minimal response
type PolemicaResponse struct {
	Body       []byte
	StatusCode int
}

func NewPolemicaRequester(cfg *pparser.MafiaBotPparserConfig) *PolemicaRequester {
	return &PolemicaRequester{
		Url:               cfg.PolemicaHost,
		Csrf:              cfg.CSRF,
		CsrfInitialCookie: cfg.CSRFCookie,
	}
}

func (p *PolemicaRequester) SetAuthCookie(header http.Header) {
	for _, newRawCookie := range header["Set-Cookie"] {
		newCookie := strings.Split(newRawCookie, "=")
		newCookieName := newCookie[0]
		newCookieBody := newCookie[1]
		for idx, cookie := range p.AuthCookie {
			if cookie.Name == newCookieName {
				p.AuthCookie[idx] = p.AuthCookie[len(p.AuthCookie)-1]
				p.AuthCookie = p.AuthCookie[:len(p.AuthCookie)-1]
				break
			}
		}
		newCookieValue := strings.Split(newCookieBody, ";")[0]
		p.AuthCookie = append(p.AuthCookie, &http.Cookie{Name: newCookieName, Value: newCookieValue})
	}
}

type QueryParams struct {
	Param string
	Value string
}

func (p *PolemicaRequester) Request(method, route string, body io.Reader, queryParams []*QueryParams) (*PolemicaResponse, error) {
	// TODO Consider refactor to kwarg patter like here https://levelup.gitconnected.com/optional-function-parameter-pattern-in-golang-c1acc829307b
	url := p.Url + route
	return p.PolemicaRequest(method, url, body, &client, queryParams)
}

func (p *PolemicaRequester) PolemicaRequest(method, url string, body io.Reader, client HttpClientInterface, queryParams []*QueryParams) (*PolemicaResponse, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return &PolemicaResponse{}, &MafiaBotPolemicaParserNewRequestError{Detail: err.Error()}
	}
	// Add Query params to request
	if queryParams != nil {
		q := req.URL.Query()
		for _, queryParam := range queryParams {
			q.Add(queryParam.Param, queryParam.Value)
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-Type", "application/json")
	if len(p.AuthCookie) == 0 {
		req.Header.Set("x-csrf-token", p.Csrf)
		req.AddCookie(&http.Cookie{
			Name:   "_csrf-maf11front",
			Value:  p.CsrfInitialCookie,
			MaxAge: 86400,
		})
	} else {
		for _, c := range p.AuthCookie {
			req.AddCookie(c)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return &PolemicaResponse{}, &MafiaBotPolemicaParserRequestConnectionError{Detail: err.Error()}
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return &PolemicaResponse{}, &MafiaBotPolemicaParserResponseBodyParsingError{Detail: err.Error()}
	}

	if resp.StatusCode > 299 {
		return &PolemicaResponse{}, &MafiaBotPolemicaParserServerResponseError{ResponseCode: resp.StatusCode, ResponseBody: string(b)}
	}

	p.SetAuthCookie(resp.Header)

	return &PolemicaResponse{b, resp.StatusCode}, nil
}
