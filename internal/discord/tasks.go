package discord

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/discord"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"log"
	"time"
)

type SendGameStatisticTaskOptions struct {
	TickerDelay int
}

type SendGameStatisticTaskOptionsHandler func(o *SendGameStatisticTaskOptions)

func SetStatisticTaskTickerDelay(delay int) func(o *SendGameStatisticTaskOptions) {
	return func(o *SendGameStatisticTaskOptions) {
		o.TickerDelay = delay
	}
}

func SendGameStatisticTask(mdb db.MafiaDBInterface, discord MafiaBotInterfaces, discordConfig discord.MafiaBotDiscordConfig, quit chan bool, opts ...SendGameStatisticTaskOptionsHandler) {
	options := &SendGameStatisticTaskOptions{
		TickerDelay: 600000, // 10 minutes in millisecond
	}
	for _, opt := range opts {
		opt(options)
	}
	ticker := time.NewTicker(time.Duration(options.TickerDelay) * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				dailyStat, err := mdb.GetDailyStatistic(discordConfig.DaySwitchHour)
				if err != nil {
					log.Printf("Error in GetDailyStatistic func: %s", err.Error())
					break
				}
				err = SendStatistics(mdb, discord, dailyStat, discordConfig.BotStatisticChannel)
				if err != nil {
					log.Printf("Error while sending stats to discord: %s", err.Error())
				}
			case <-quit:
				return
			}
		}
	}()
}
