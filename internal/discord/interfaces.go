package discord

import "github.com/bwmarrin/discordgo"

type MafiaBotInterfaces interface {
	GuildChannels(string) ([]*discordgo.Channel, error)
	ChannelMessageSend(string, string) (*discordgo.Message, error)
}
