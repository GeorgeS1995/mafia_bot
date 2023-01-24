package cfg

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/db"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/pparser"
)

type MafiaBotGeneralConfig struct {
	DB      *db.MafiaBotDBConfig
	Discord *discord.MafiaBotDiscordConfig
	Pparser *pparser.MafiaBotPparserConfig
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
	pparserConfig, err := pparser.NewMafiaBotPparserConfig()
	if err != nil {
		return &MafiaBotGeneralConfig{}, err
	}
	return &MafiaBotGeneralConfig{
		dbConfig,
		discordConfing,
		pparserConfig,
	}, nil
}
