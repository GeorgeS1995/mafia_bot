package cfg

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/db"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
)

type MafiaBotGeneralConfig struct {
	DB      *db.MafiaBotDBConfig
	Discord *discord.MafiaBotDiscordConfig
}

func NewMafiaBotGeneralConfig() (*MafiaBotGeneralConfig, error) {
	dbConfig, err := db.NewMafiaBotDBConfig()
	if err != nil {
		return &MafiaBotGeneralConfig{}, err
	}
	discordConfing, err := discord.NewMafiaBotDiscordConfig()
	if err != nil {
		return &MafiaBotGeneralConfig{}, err
	}
	return &MafiaBotGeneralConfig{
		dbConfig,
		discordConfing,
	}, nil
}
