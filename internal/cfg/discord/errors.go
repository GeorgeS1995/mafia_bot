package discord

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
)

type MafiaBotConfigDaySwitchHourError struct {
	common.MafiaBotConfigError
	DaySwitchHour int
}

func (e *MafiaBotConfigDaySwitchHourError) Error() string {
	return fmt.Sprintf("%v: Day switch hour must be between 0 and 23, actual value: %d", e.GetISOFormat(), e.DaySwitchHour)
}
