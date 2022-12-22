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

type MafiaBotParseError struct {
	MafiaBotConfigError
	ParsedAttr string
}

func (e *MafiaBotParseError) Error() string {
	return fmt.Sprintf("%v: DB env %s is requierd parameter", e.GetISOFormat(), e.ParsedAttr)
}
