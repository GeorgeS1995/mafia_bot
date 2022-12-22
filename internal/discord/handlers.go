package discord

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
	"log"
)

// TODO refactor to goroutin
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
