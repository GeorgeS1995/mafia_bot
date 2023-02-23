package discord

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/GeorgeS1995/mafia_bot/test"
	"os"
	"testing"
)

func TestGetTokenOK(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_DISCORD_TOKEN")
	}()
	token := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
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
