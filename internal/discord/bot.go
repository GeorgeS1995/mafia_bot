package discord

import (
	"github.com/bwmarrin/discordgo"
)

type MafiaDiscordBot struct {
	Bot      *discordgo.Session
	handlers []interface{}
}

func NewMafiaDiscordBot(token string) (*MafiaDiscordBot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return &MafiaDiscordBot{}, err
	}
	handlers := []interface{}{
		greetings,
	}
	return &MafiaDiscordBot{Bot: session, handlers: handlers}, nil
}

func (b *MafiaDiscordBot) Init() error {
	for _, handler := range b.handlers {
		b.Bot.AddHandler(handler)
	}
	err := b.Bot.Open()
	if err != nil {
		return err
	}
	return nil
}
