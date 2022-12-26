package discord

import (
	"errors"
	"github.com/GeorgeS1995/mafia_bot/internal/discord"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"os"
	"testing"
)

func TestGreetingsGuildChannelsError(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_STATUS_CHANNELS")
		_ = os.Unsetenv("MAFIA_BOT_DISCORD_TOKEN")

	}()
	token := test.RandStringRunes(3)
	statusChannel := test.RandStringRunes(3) + ","
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
	_ = os.Setenv("MAFIA_BOT_STATUS_CHANNELS", statusChannel)
	guilds := []*discordgo.Guild{{ID: test.RandStringRunes(3)}}
	guildChannelsError := errors.New(test.RandStringRunes(3))
	expectedError := &discord.MafiaBotGuildChannelsError{
		GuildId: guilds[0].ID,
		Err:     guildChannelsError,
	}
	ctrl := gomock.NewController(t)
	m := NewMockMafiaBotInterfaces(ctrl)
	m.
		EXPECT().
		GuildChannels(gomock.Any()).
		Return([]*discordgo.Channel{}, guildChannelsError).
		Times(1)

	err := discord.Greetings(m, guilds)

	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("Nil error or not expected error.\n Expected: %v\n Passed: %v\n", expectedError, err)
	}
}

func TestGreetingsChannelMessageSendError(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_STATUS_CHANNELS")
		_ = os.Unsetenv("MAFIA_BOT_DISCORD_TOKEN")

	}()
	token := test.RandStringRunes(3)
	statusChannel := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
	_ = os.Setenv("MAFIA_BOT_STATUS_CHANNELS", statusChannel+",")
	channelsError := errors.New(test.RandStringRunes(3))
	expectedError := &discord.MafiaBotChannelMSGSendError{
		ChannelId: statusChannel,
		Err:       channelsError,
	}
	ctrl := gomock.NewController(t)
	m := NewMockMafiaBotInterfaces(ctrl)
	m.
		EXPECT().
		GuildChannels(gomock.Any()).
		Return([]*discordgo.Channel{{ID: statusChannel}}, nil).
		Times(1)
	m.
		EXPECT().
		ChannelMessageSend(statusChannel, discord.BotOnline).
		Return(&discordgo.Message{}, channelsError).
		Times(1)

	err := discord.Greetings(m, []*discordgo.Guild{{}})

	if err == nil || err.Error() != expectedError.Error() {
		t.Fatalf("Nil error or not expected error.\n Expected: %v\n Passed: %v\n", expectedError, err)
	}
}

func TestGreetingsEmptyBotStatusChannels(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_DISCORD_TOKEN")

	}()
	token := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
	ctrl := gomock.NewController(t)
	m := NewMockMafiaBotInterfaces(ctrl)
	m.
		EXPECT().
		GuildChannels(gomock.Any()).
		Times(0)
	m.
		EXPECT().
		ChannelMessageSend(gomock.Any(), discord.BotOnline).
		Times(0)

	err := discord.Greetings(m, []*discordgo.Guild{{}})

	if err != nil {
		t.Fatalf("Not nil error.\n Passed: %v\n", err)
	}
}
