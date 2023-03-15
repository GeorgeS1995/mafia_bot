package discord

import (
	"errors"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/discord"
	"github.com/GeorgeS1995/mafia_bot/test"
	mockdb "github.com/GeorgeS1995/mafia_bot/test/internal"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestGreetingsGuildChannelsError(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_STATUS_CHANNELS")
		_ = os.Unsetenv("MAFIA_BOT_DISCORD_TOKEN")
		_ = os.Unsetenv("MAFIA_BOT_STATISTIC_CHANNEL")
	}()
	token := test.RandStringRunes(3)
	statusChannel := test.RandStringRunes(3) + ","
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
	_ = os.Setenv("MAFIA_BOT_STATUS_CHANNELS", statusChannel)
	_ = os.Setenv("MAFIA_BOT_STATISTIC_CHANNEL", test.RandStringRunes(3))
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
		_ = os.Unsetenv("MAFIA_BOT_STATISTIC_CHANNEL")

	}()
	token := test.RandStringRunes(3)
	statusChannel := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DISCORD_TOKEN", token)
	_ = os.Setenv("MAFIA_BOT_STATUS_CHANNELS", statusChannel+",")
	_ = os.Setenv("MAFIA_BOT_STATISTIC_CHANNEL", test.RandStringRunes(3))
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

func TestSendStatisticsOK(t *testing.T) {
	playedGames := []uuid.UUID{uuid.New()}
	dailyStat := []*db.DailyStatistic{{
		PlayedGamesID: playedGames,
	}}
	statisticChannel := test.RandStringRunes(3)
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockDiscordBot := NewMockMafiaBotInterfaces(ctrl)
	gameStatisticMsg := `Статистика за день 01-01-0001
Количество сыгранных игр: 0
Победы мафии: 0
Победы мирных: 0

MVP вечера  с суммой очков 0.00
`
	mockDiscordBot.EXPECT().ChannelMessageSend(statisticChannel, gameStatisticMsg).Return(&discordgo.Message{}, nil)
	mockMafiaDB.EXPECT().MarkGamesAsSent(playedGames).Return(nil)

	err := discord.SendStatistics(mockMafiaDB, mockDiscordBot, dailyStat, statisticChannel)

	if err != nil {
		t.Fatalf("Not nil error.\n Passed: %v\n", err)
	}
}

func TestSendStatisticsChannelMessageSendError(t *testing.T) {
	playedGames := []uuid.UUID{uuid.New()}
	dailyStat := []*db.DailyStatistic{{
		PlayedGamesID: playedGames,
	}}
	statisticChannel := test.RandStringRunes(3)
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockDiscordBot := NewMockMafiaBotInterfaces(ctrl)
	mockDiscordBot.EXPECT().ChannelMessageSend(statisticChannel, gomock.Any()).Return(&discordgo.Message{}, errors.New(""))
	mockMafiaDB.EXPECT().MarkGamesAsSent(playedGames).Return(nil).Times(0)

	err := discord.SendStatistics(mockMafiaDB, mockDiscordBot, dailyStat, statisticChannel)

	if _, ok := err.(*discord.MafiaBotSendStatisticsError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestSendStatisticsMarkGamesAsSentError(t *testing.T) {
	playedGames := []uuid.UUID{uuid.New()}
	dailyStat := []*db.DailyStatistic{{
		PlayedGamesID: playedGames,
	}}
	statisticChannel := test.RandStringRunes(3)
	ctrl := gomock.NewController(t)
	mockMafiaDB := mockdb.NewMockMafiaDBInterface(ctrl)
	mockDiscordBot := NewMockMafiaBotInterfaces(ctrl)
	mockDiscordBot.EXPECT().ChannelMessageSend(statisticChannel, gomock.Any()).Return(&discordgo.Message{}, nil)
	mockMafiaDB.EXPECT().MarkGamesAsSent(playedGames).Return(errors.New(""))

	err := discord.SendStatistics(mockMafiaDB, mockDiscordBot, dailyStat, statisticChannel)

	if _, ok := err.(*discord.MafiaBotSendStatisticsError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}
