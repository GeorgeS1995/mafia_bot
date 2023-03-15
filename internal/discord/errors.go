package discord

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal"
)

type MafiaBotDiscordApiError struct {
	internal.MafiaBotError
}

func (e *MafiaBotDiscordApiError) Error() string {
	return fmt.Sprintf("%v: Discord api answer error", e.GetISOFormat())
}

type MafiaBotGuildChannelsError struct {
	internal.MafiaBotError
	GuildId string
	Err     error
}

func (e *MafiaBotGuildChannelsError) Error() string {
	return fmt.Sprintf("%v: Can't get channels from guildID: %s, original error: %s", e.GetISOFormat(), e.GuildId, e.Err.Error())
}

type MafiaBotChannelMSGSendError struct {
	internal.MafiaBotError
	ChannelId string
	Err       error
}

func (e *MafiaBotChannelMSGSendError) Error() string {
	return fmt.Sprintf("%v: Can't send msg to the channel id: %s, original error: %s", e.GetISOFormat(), e.ChannelId, e.Err.Error())
}

type MafiaBotSendStatisticsError struct {
	internal.MafiaBotError
	SentErrorsList []dailyStatisticErrorData
}

func (e *MafiaBotSendStatisticsError) Error() string {
	formatedStatisticErrors := ""
	for _, err := range e.SentErrorsList {
		formatedStatisticErrors += fmt.Sprintf("Date: %s, From: %s, Description: %s\n", err.Day, err.From, err.Detail)
	}
	return fmt.Sprintf("%v: This games was not send to discord: \n", e.GetISOFormat())
}
