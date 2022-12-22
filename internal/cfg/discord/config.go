package discord

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"os"
	"strings"
)

type MafiaBotDiscordConfig struct {
	Token             string
	BotStatusChannels []string
}

func NewMafiaBotDiscordConfig() (*MafiaBotDiscordConfig, error) {
	discordConfig := &MafiaBotDiscordConfig{}
	token, err := discordConfig.GetToken()
	if err != nil {
		return discordConfig, err
	}
	botStatusChannels := discordConfig.GetBotStatusChannels()
	discordConfig.Token = token
	discordConfig.BotStatusChannels = botStatusChannels
	return discordConfig, nil
}

func (c *MafiaBotDiscordConfig) GetToken() (token string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "DISCORD_TOKEN")
	token = os.Getenv(env_name)
	if token == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return token, err
}

func (c *MafiaBotDiscordConfig) GetBotStatusChannels() []string {
	env_name := fmt.Sprintf(common.ConfPrefix, "STATUS_CHANNELS")
	channelsString := os.Getenv(env_name)
	if channelsString == "" {
		return []string{}
	}
	return strings.Split(channelsString, ",")
}