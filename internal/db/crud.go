package db

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func (b *MafiaDB) UpdateOrCreateUser(user *User) (*User, error) {
	foundedUser := &User{}

	result := b.Db.Where(&User{PolemicaId: user.PolemicaId}).First(foundedUser)
	if result.Error != nil && !(errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		return user, &MafiaBotUpdateOrCreateUserGetError{Detail: result.Error.Error()}
	}
	user.ID = foundedUser.ID
	user.CreatedAt = foundedUser.CreatedAt
	result = b.Db.Save(user)
	if result.Error != nil {
		return user, &MafiaBotUpdateOrCreateUserSaveError{Detail: result.Error.Error()}
	}
	return user, nil
}

func (b *MafiaDB) GetLastGame() (*Game, error) {
	game := &Game{}
	result := b.Db.Order("started_at desc").First(game)
	var err error
	if result.Error != nil && !(errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		err = &MafiaBotGetLastGameDriverError{Detail: result.Error.Error()}
	} else if result.Error != nil && (errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		err = &MafiaBotGetLastGameEmptyDBError{}
	}
	return game, err
}

func (b *MafiaDB) Create(value interface{}) error {
	result := b.Db.Create(value)
	return result.Error
}

// Transaction with respect b.Db attribute
func (b *MafiaDB) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	originalDB := b.Db
	defer func() {
		b.Db = originalDB
	}()
	fcDecorator := func(tx *gorm.DB) error {
		b.Db = tx
		return fc(tx)
	}
	return b.Db.Transaction(fcDecorator, opts...)
}

// GetDailyStatistic Assume that polemica nickname is unique
func (b *MafiaDB) GetDailyStatistic(DaySwitchHour int) ([]*DailyStatistic, error) {
	timeNow := time.Now()
	rows, err := b.Db.Raw(`select 
								  g.id, 
								  g.winner,
								  TIMESTAMP 'epoch' + INTERVAL '1 second' * trunc(
									(
									  extract(
										'epoch' 
										FROM 
										  g.started_at
									  ) / 86400
									)
								  ) * 86400 AS clock, 
								  u.polemica_nick_name, 
								  sum(pg.points), 
								  max(
									sum(pg.points)
								  ) OVER (
									PARTITION BY u.polemica_nick_name, 
									TIMESTAMP 'epoch' + INTERVAL '1 second' * trunc(
									  (
										extract(
										  'epoch' 
										  FROM 
											g.started_at
										) / 86400
									  )
									) * 86400
								  ) as max_points 
								from 
								  games g 
								  join player_games pg on pg.game_id = g.id 
								  join users u on u.id = pg.user_id 
								where 
								  g.is_sended_to_discord = false 
								  and g.started_at <= ?
								group by 
								  g.id, 
								  u.polemica_nick_name, 
								  clock
								order by g.started_at
								`, time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), DaySwitchHour, 0, 0, 0, timeNow.Location())).Rows()
	if err != nil {
		return nil, &MafiaBotGetDailyStatisticQueryError{Detail: err.Error()}
	}
	defer rows.Close()
	daysStatistics := []*DailyStatistic{}
	for rows.Next() {
		var id uuid.UUID
		var winner GameResult
		var clock time.Time
		var polemicka_nick_name string
		var sum float32
		var max_points float32
		err = rows.Scan(&id, &winner, &clock, &polemicka_nick_name, &sum, &max_points)
		if err != nil {
			return nil, &MafiaBotGetDailyStatisticRowScanError{Detail: err.Error()}
		}
		daysStatisticsLen := len(daysStatistics)
		roundedClock := time.Date(clock.Year(), clock.Month(), clock.Day(), 0, 0, 0, 0, clock.Location())
		if daysStatisticsLen == 0 || daysStatistics[daysStatisticsLen-1].Date.Before(roundedClock) {
			mafiaWins := 0
			cityWins := 0
			if winner == CityWin {
				cityWins++
			}
			if winner == MafiaWin {
				mafiaWins++
			}
			daysStatistics = append(daysStatistics, &DailyStatistic{
				Date:      time.Date(clock.Year(), clock.Month(), clock.Day(), 0, 0, 0, 0, clock.Location()),
				GameCount: 1,
				MafiaWins: mafiaWins,
				CityWins:  cityWins,
				MVP: MVP{
					polemicka_nick_name,
					max_points,
				},
				PlayedGamesID: []uuid.UUID{id},
			})
		} else {
			lastDayStatistic := daysStatistics[daysStatisticsLen-1]
			if lastDayStatistic.MVP.Score < max_points {
				lastDayStatistic.MVP = MVP{
					polemicka_nick_name,
					max_points,
				}
			}
			IsGameCounted := false
			for _, gameID := range lastDayStatistic.PlayedGamesID {
				if gameID == id {
					IsGameCounted = true
					break
				}
			}
			if !(IsGameCounted) {
				if winner == CityWin {
					lastDayStatistic.CityWins++
				}
				if winner == MafiaWin {
					lastDayStatistic.MafiaWins++
				}
				lastDayStatistic.GameCount++
				lastDayStatistic.PlayedGamesID = append(lastDayStatistic.PlayedGamesID, id)
			}
		}
	}
	return daysStatistics, nil
}

func (b *MafiaDB) MarkGamesAsSent(gameIDS []uuid.UUID) (err error) {
	games := []*Game{}
	result := b.Db.Find(&games, gameIDS)
	if result.Error != nil {
		return &MafiaBotMarkGamesAsSentFindError{
			Games:  gameIDS,
			Detail: result.Error.Error(),
		}
	}
	err = b.Transaction(func(tx *gorm.DB) error {
		for _, g := range games {
			g.IsSendedToDiscord = true
			result = tx.Save(g)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
	if err != nil {
		err = &MafiaBotMarkGamesAsSentTransactionError{
			Games:  gameIDS,
			Detail: err.Error(),
		}
	}
	return err
}
