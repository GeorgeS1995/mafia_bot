package internal

import (
	"fmt"
	"time"
)

type MafiaBotError struct{}

func (e *MafiaBotError) GetISOFormat() string {
	return time.Now().Format(time.RFC3339)
}

func (e *MafiaBotError) Error() string {
	return fmt.Sprintf("%v: General mafia bot error", e.GetISOFormat())
}
