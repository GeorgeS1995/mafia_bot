package discord

import (
	"errors"
	cfg "github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/discord"
	"github.com/GeorgeS1995/mafia_bot/test"
	test_common "github.com/GeorgeS1995/mafia_bot/test/internal"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestSendGameStatisticTaskOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	statisticChannel := test.RandStringRunes(3)
	discordConfig := cfg.MafiaBotDiscordConfig{BotStatisticChannel: statisticChannel}
	mockMafiaDB := test_common.NewMockMafiaDBInterface(ctrl)
	mockMafiaDB.EXPECT().GetDailyStatistic(discordConfig.DaySwitchHour).Return([]*db.DailyStatistic{{}}, nil).Times(2)
	mockMafiaDB.EXPECT().MarkGamesAsSent(gomock.Any()).Return(nil).Times(2)
	mockDiscordBot := NewMockMafiaBotInterfaces(ctrl)
	mockDiscordBot.EXPECT().ChannelMessageSend(statisticChannel, gomock.Any()).Return(&discordgo.Message{}, nil).Times(2)
	SendGameStatisticTaskQuiteChan := make(chan bool)

	discord.SendGameStatisticTask(mockMafiaDB, mockDiscordBot, discordConfig, SendGameStatisticTaskQuiteChan, discord.SetStatisticTaskTickerDelay(test_common.TestTick))

	time.Sleep(time.Duration(test_common.TestTick*2) * time.Millisecond)
	SendGameStatisticTaskQuiteChan <- true
}

func TestSendGameStatisticTaskDailyStatError(t *testing.T) {
	ctrl := gomock.NewController(t)
	statisticChannel := test.RandStringRunes(3)
	discordConfig := cfg.MafiaBotDiscordConfig{BotStatisticChannel: statisticChannel}
	mockMafiaDB := test_common.NewMockMafiaDBInterface(ctrl)
	GetDailyStatisticFirstCall := true
	mockMafiaDB.EXPECT().GetDailyStatistic(discordConfig.DaySwitchHour).DoAndReturn(func(DaySwitchHour int) ([]*db.DailyStatistic, error) {
		if GetDailyStatisticFirstCall {
			GetDailyStatisticFirstCall = false
			return []*db.DailyStatistic{{}}, nil
		} else {
			return []*db.DailyStatistic{{}}, errors.New("")
		}
	}).Times(2)
	mockMafiaDB.EXPECT().MarkGamesAsSent(gomock.Any()).Return(nil).Times(1)
	mockDiscordBot := NewMockMafiaBotInterfaces(ctrl)
	mockDiscordBot.EXPECT().ChannelMessageSend(statisticChannel, gomock.Any()).Return(&discordgo.Message{}, nil).Times(1)
	SendGameStatisticTaskQuiteChan := make(chan bool)

	discord.SendGameStatisticTask(mockMafiaDB, mockDiscordBot, discordConfig, SendGameStatisticTaskQuiteChan, discord.SetStatisticTaskTickerDelay(test_common.TestTick))

	time.Sleep(time.Duration(test_common.TestTick*2) * time.Millisecond)
	SendGameStatisticTaskQuiteChan <- true
}
