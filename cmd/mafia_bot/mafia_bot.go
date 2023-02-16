package main

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg"
	db2 "github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/internal/discord"
	"github.com/GeorgeS1995/mafia_bot/internal/pparser"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := cfg.NewMafiaBotGeneralConfig()
	if err != nil {
		log.Fatal("Can't parse config, ", err)
		return
	}
	db, err := db2.NewMafiaDB(*config.DB)
	if err != nil {
		log.Fatal("Can't create db connection, ", err)
		return
	}
	err = db.Init()
	if err != nil {
		log.Fatal("Can't init db migration, ", err)
		return
	}
	dg, err := discord.NewMafiaDiscordBot(config.Discord.Token)
	if err != nil {
		log.Fatal("Error creating Discord session, ", err)
		return
	}
	err = dg.Init()
	if err != nil {
		log.Fatal("Can't create discord client session, ", err)
		return
	}
	pc := pparser.NewPolemicaApiClient(config.Pparser)
	err = pc.Login(config.Pparser.Login, config.Pparser.Password)
	if err != nil {
		log.Fatal("Can't login on polemica site, ", err)
		return
	}
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Bot.Close()
}
