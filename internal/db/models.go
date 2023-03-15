package db

import (
	"github.com/google/uuid"
	"time"
)

type MafiaBaseModel struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Game struct {
	MafiaBaseModel
	PolemicaId        string
	StartedAt         time.Time
	IsSendedToDiscord bool
	Winner            GameResult `gorm:"type:game_result"`
	Players           []User     `gorm:"many2many:player_games;"`
}

type PlayerGame struct {
	UserID uuid.UUID `gorm:"index:user_game_id,unique;type:uuid;"`
	GameID uuid.UUID `gorm:"index:user_game_id,unique;type:uuid;"`
	Points float32
}

// User TODO move to the separated table info about thirdPartyServices
type User struct {
	MafiaBaseModel
	PolemicaNickName string
	PolemicaId       string
}

type MVP struct {
	NickName string
	Score    float32
}
type DailyStatistic struct {
	Date          time.Time
	GameCount     int
	MafiaWins     int
	CityWins      int
	MVP           MVP
	PlayedGamesID []uuid.UUID
}
