package common

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal"
)

type MafiaBotConfigError struct {
	internal.MafiaBotError
}

func (e *MafiaBotConfigError) Error() string {
	return fmt.Sprintf("%v: Config parse mafia bot error", e.GetISOFormat())
}

type MafiaBotParseMissingRequiredParamError struct {
	MafiaBotConfigError
	ParsedAttr string
}

func (e *MafiaBotParseMissingRequiredParamError) Error() string {
	return fmt.Sprintf("%v: Env %s is requierd parameter", e.GetISOFormat(), e.ParsedAttr)
}

type MafiaBotParseTypeError struct {
	MafiaBotConfigError
	ParsedAttr string
}

func (e *MafiaBotParseTypeError) Error() string {
	return fmt.Sprintf("%v: Wrong type of env: %s", e.GetISOFormat(), e.ParsedAttr)
}
