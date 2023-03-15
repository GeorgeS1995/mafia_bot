package discord

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
	"log"
	"time"
)

// Greetings TODO refactor to goroutin
func Greetings(s MafiaBotInterfaces, guilds []*discordgo.Guild) error {
	// Is there a flexible way to pass config
	discordConfig, _ := discord.NewMafiaBotDiscordConfig()
	if len(discordConfig.BotStatusChannels) == 0 {
		return nil
	}

	for _, guild := range guilds {
		channelsList, err := s.GuildChannels(guild.ID)
		if err != nil {
			return &MafiaBotGuildChannelsError{
				GuildId: guild.ID,
				Err:     err,
			}
		}
		for _, channel := range channelsList {
			if slices.Contains(discordConfig.BotStatusChannels, channel.ID) {
				_, err = s.ChannelMessageSend(channel.ID, BotOnline)
				if err != nil {
					return &MafiaBotChannelMSGSendError{
						ChannelId: channel.ID,
						Err:       err,
					}
				}
			}
		}
	}
	return nil
}

func greetings(s *discordgo.Session, m *discordgo.Connect) {
	guilds := s.State.Guilds
	err := Greetings(s, guilds)
	if err != nil {
		log.Println(err.Error())
	}
}

type dailyStatisticErrorData struct {
	Day    time.Time
	From   SendErrorFrom
	Detail string
}

func SendStatistics(mdb db.MafiaDBInterface, s MafiaBotInterfaces, dailyStat []*db.DailyStatistic, statisticChannel string) error {
	dailyStatisticErrors := []dailyStatisticErrorData{}
	for _, game := range dailyStat {
		_, err := s.ChannelMessageSend(statisticChannel, fmt.Sprintf(GameStatistic, game.Date.Format(internal.UserDataFormat), game.GameCount, game.MafiaWins, game.CityWins, game.MVP.NickName, game.MVP.Score))
		if err != nil {
			dailyStatisticErrors = append(dailyStatisticErrors, dailyStatisticErrorData{
				game.Date,
				Discord,
				err.Error(),
			})
			continue
		}
		err = mdb.MarkGamesAsSent(game.PlayedGamesID)
		if err != nil {
			dailyStatisticErrors = append(dailyStatisticErrors, dailyStatisticErrorData{
				game.Date,
				DB,
				err.Error(),
			})
		}
	}
	var err error
	if len(dailyStatisticErrors) > 0 {
		err = &MafiaBotSendStatisticsError{SentErrorsList: dailyStatisticErrors}
	}
	return err
}
