package pparser

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal"
)

type MafiaBotPolemicaParserError struct {
	internal.MafiaBotError
}

type MafiaBotPolemicaParserNewRequestError struct {
	MafiaBotPolemicaParserError
	Detail string
}

func (e *MafiaBotPolemicaParserNewRequestError) Error() string {
	return fmt.Sprintf("%v: Can't create new request object: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotPolemicaParserRequestConnectionError struct {
	MafiaBotPolemicaParserError
	Detail string
}

func (e *MafiaBotPolemicaParserRequestConnectionError) Error() string {
	return fmt.Sprintf("%v: Polemica request connection error: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotPolemicaParserServerResponseError struct {
	MafiaBotPolemicaParserError
	ResponseCode int
	ResponseBody string
}

func (e *MafiaBotPolemicaParserServerResponseError) Error() string {
	return fmt.Sprintf("%v: Not 200x response.\n Code: %s\n Body: %s\n", e.GetISOFormat(), e.ResponseCode, e.ResponseBody)
}

type MafiaBotPolemicaParserResponseBodyParsingError struct {
	MafiaBotPolemicaParserError
	Detail string
}

func (e *MafiaBotPolemicaParserResponseBodyParsingError) Error() string {
	return fmt.Sprintf("%v: Can't parse response body: %s", e.GetISOFormat(), e.Detail)
}

type MafiaBotPolemicaParserParseGameResponseError struct {
	MafiaBotPolemicaParserError
	Detail string
	GameID string
}

func (e *MafiaBotPolemicaParserParseGameResponseError) Error() string {
	return fmt.Sprintf("%v: Can't get game %s info, detail: %s", e.GetISOFormat(), e.GameID, e.Detail)
}

type MafiaBotPolemicaParserParseGameUnmarshalError struct {
	MafiaBotPolemicaParserError
	Detail string
	GameID string
}

func (e *MafiaBotPolemicaParserParseGameUnmarshalError) Error() string {
	return fmt.Sprintf("%v: Can't unmarshal game %s info, detail: %s", e.GetISOFormat(), e.GameID, e.Detail)
}

type MafiaBotPolemicaParserParseGamesHistoryResponseError struct {
	MafiaBotPolemicaParserError
	Detail     string
	QueryParam []*QueryParams
}

func (e *MafiaBotPolemicaParserParseGamesHistoryResponseError) Error() string {
	queryString := ""
	for _, q := range e.QueryParam {
		queryString += fmt.Sprintf("param: %s, value: %s", q.Param, q.Value)
	}
	return fmt.Sprintf("%v: Can't get game history info.\n detail: %s\n query: %s", e.GetISOFormat(), e.Detail, queryString)
}

type MafiaBotPolemicaParserParseGamesHistoryUnmarshalError struct {
	MafiaBotPolemicaParserError
	Detail string
}

func (e *MafiaBotPolemicaParserParseGamesHistoryUnmarshalError) Error() string {
	return fmt.Sprintf("%v: Can't unmarshal game history info, detail: %s", e.GetISOFormat(), e.Detail)
}
