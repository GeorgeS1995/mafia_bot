package discord

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"os"
	"strconv"
	"strings"
)

type MafiaBotDiscordConfig struct {
	Token               string
	BotStatusChannels   []string
	BotStatisticChannel string
	DaySwitchHour       int
}

func NewMafiaBotDiscordConfig() (*MafiaBotDiscordConfig, error) {
	discordConfig := &MafiaBotDiscordConfig{}
	token, err := discordConfig.GetToken()
	if err != nil {
		return discordConfig, err
	}
	botStatusChannels := discordConfig.GetBotStatusChannels()
	statisticChannel, err := discordConfig.GetBotStatisticChannel()
	if err != nil {
		return discordConfig, err
	}
	daySwitchHour, err := discordConfig.GetDaySwitchHour()
	if err != nil {
		return discordConfig, err
	}
	discordConfig.Token = token
	discordConfig.BotStatusChannels = botStatusChannels
	discordConfig.BotStatisticChannel = statisticChannel
	discordConfig.DaySwitchHour = daySwitchHour
	return discordConfig, nil
}

func (c *MafiaBotDiscordConfig) GetToken() (token string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "DISCORD_TOKEN")
	token = os.Getenv(envName)
	if token == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return token, err
}

func (c *MafiaBotDiscordConfig) GetBotStatusChannels() []string {
	envName := fmt.Sprintf(common.ConfPrefix, "STATUS_CHANNELS")
	channelsString := os.Getenv(envName)
	if channelsString == "" {
		return []string{}
	}
	return strings.Split(channelsString, ",")
}

func (c *MafiaBotDiscordConfig) GetBotStatisticChannel() (channel string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "STATISTIC_CHANNEL")
	channel = os.Getenv(envName)
	if channel == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return channel, err
}

func (c *MafiaBotDiscordConfig) GetDaySwitchHour() (daySwitchHour int, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "DAY_SWITCH_HOUR")
	daySwitchHourStr := os.Getenv(envName)
	if daySwitchHourStr != "" {
		daySwitchHour, err = strconv.Atoi(daySwitchHourStr)
		if err != nil {
			err = &common.MafiaBotParseTypeError{ParsedAttr: envName}
		}
		if daySwitchHour < 0 || daySwitchHour > 23 {
			err = &MafiaBotConfigDaySwitchHourError{DaySwitchHour: daySwitchHour}
		}
	}
	return daySwitchHour, err
}
