package discord

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestGetTokenOK(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_DISCORD_TOKEN")
		_ = os.Unsetenv("MAFIA_BOT_STATISTIC_CHANNEL")
	}()
	token := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
	_ = os.Setenv("MAFIA_BOT_STATISTIC_CHANNEL", test.RandStringRunes(3))
	discordConfig, err := discord.NewMafiaBotDiscordConfig()

	if err != nil {
		t.Fatalf("Unexpected error while parsing discord config: %v", err)
	}
	if discordConfig.Token != token {
		t.Fatalf("Token %v is not equal expected.", discordConfig.Token)
	}
}

func TestGetTokenNotSet(t *testing.T) {
	discordConfig, err := discord.NewMafiaBotDiscordConfig()

	expectedErr := common.MafiaBotParseMissingRequiredParamError{ParsedAttr: "MAFIA_BOT_DISCORD_TOKEN"}
	expectedErrorMsg := expectedErr.Error()
	if err == nil || err.Error() != expectedErrorMsg {
		t.Fatalf("Unexpected error while parsing discord config: %v", err)
	}
	if discordConfig.Token != "" {
		t.Fatalf("Token %v is not equal expected.", discordConfig.Token)
	}
}

func TestGetDaySwitchHourOK(t *testing.T) {
	defer os.Unsetenv("MAFIA_BOT_DAY_SWITCH_HOUR")
	discordConfig := &discord.MafiaBotDiscordConfig{}
	for i := 0; i < 24; i++ {
		t.Run(fmt.Sprintf("DAY_SWITCH_HOUR = %d", i), func(t *testing.T) {
			_ = os.Setenv("MAFIA_BOT_DAY_SWITCH_HOUR", strconv.Itoa(i))

			daySwitchHour, err := discordConfig.GetDaySwitchHour()

			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}
			assert.Equal(t, i, daySwitchHour)
		})
	}
}

func TestGetDaySwitchHourDefaultOK(t *testing.T) {
	_ = os.Unsetenv("MAFIA_BOT_DAY_SWITCH_HOUR")
	discordConfig := &discord.MafiaBotDiscordConfig{}

	daySwitchHour, err := discordConfig.GetDaySwitchHour()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	assert.Equal(t, 0, daySwitchHour)
}

func TestGetDaySwitchHourParseTypeError(t *testing.T) {
	defer os.Unsetenv("MAFIA_BOT_DAY_SWITCH_HOUR")
	discordConfig := &discord.MafiaBotDiscordConfig{}
	_ = os.Setenv("MAFIA_BOT_DAY_SWITCH_HOUR", test.RandStringRunes(3))

	_, err := discordConfig.GetDaySwitchHour()

	if _, ok := err.(*common.MafiaBotParseTypeError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestGetDaySwitchHourTimeRangeError(t *testing.T) {
	defer os.Unsetenv("MAFIA_BOT_DAY_SWITCH_HOUR")
	discordConfig := &discord.MafiaBotDiscordConfig{}
	for _, hour := range [2]int{-1, 24} {
		t.Run(fmt.Sprintf("DAY_SWITCH_HOUR = %d", hour), func(t *testing.T) {
			_ = os.Setenv("MAFIA_BOT_DAY_SWITCH_HOUR", strconv.Itoa(hour))

			_, err := discordConfig.GetDaySwitchHour()

			if _, ok := err.(*discord.MafiaBotConfigDaySwitchHourError); !ok {
				t.Fatalf("Wrong error type: %s", err)
			}
		})
	}
}
