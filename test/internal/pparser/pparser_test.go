package test_pparser

import (
	"bytes"
	"errors"
	"fmt"
	cfg "github.com/GeorgeS1995/mafia_bot/internal/cfg/pparser"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http"
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

func TestPolemicaParserSetUserIDError(t *testing.T) {
	pParser := pparser.NewPolemicaRequester(&cfg.MafiaBotPparserConfig{})
	method := "GET"
	url := ""
	ctrl := gomock.NewController(t)
	m := NewMockHttpClientInterface(ctrl)
	m.EXPECT().Do(gomock.Any()).Return(
		&http.Response{
			Body:       io.NopCloser(bytes.NewReader([]byte{255})),
			StatusCode: 200,
			Header:     http.Header{"Set-Cookie": []string{"%u0"}}}, nil)

	_, err := pParser.PolemicaRequest(method, url, nil, m, nil)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserSetUserIDError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestSetUserIDFromCookieOK(t *testing.T) {
	pParser := pparser.NewPolemicaRequester(&cfg.MafiaBotPparserConfig{})
	rand.Seed(time.Now().UnixNano())
	userId := rand.Intn(99999-10000) + 10000
	headers := http.Header{
		"Set-Cookie": []string{
			fmt.Sprintf("_id-maf11front=2bf4533645a4b8e8333d962ef510fa1ca620039590fff1bebdf9d438acdd3c30a:2:{i:0;s:14:\"_id-maf11front\";i:1;s:48:\"[%d,\"wqHCFqlkZzAEziaDtdENn9JH7oxR7a02\",86400]\";}; expires=Fri, 17-Feb-2023 18:47:23 GMT; Max-Age=86400; path=/; HttpOnly", userId),
			"region=694203d7630f1597b0bf7afa9f132a14091b12570991e5d855f05814064be272a:2:{i:0;s:6:\"region\";i:1;s:2:\"KZ\";}; path=/; HttpOnly",
		},
	}

	err := pParser.SetUserIDFromCookie(headers)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	assert.Equal(t, pParser.UserID, userId)
}

func TestSetUserIDFromCookieDecodeError(t *testing.T) {
	pParser := pparser.NewPolemicaRequester(&cfg.MafiaBotPparserConfig{})
	headers := http.Header{
		"Set-Cookie": []string{"%u0"},
	}

	err := pParser.SetUserIDFromCookie(headers)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserSetUserIDFromCookieDecodeError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestSetUserIDFromCookieRegexError(t *testing.T) {
	pParser := pparser.NewPolemicaRequester(&cfg.MafiaBotPparserConfig{})
	headers := http.Header{
		"Set-Cookie": []string{"_id-maf11front=2bf4"},
	}

	err := pParser.SetUserIDFromCookie(headers)

	if _, ok := err.(*pparser.MafiaBotPolemicaParserSetUserIDFromCookieRegexError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}
