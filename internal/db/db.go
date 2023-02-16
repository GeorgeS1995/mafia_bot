package db

import "fmt"

func (b *MafiaDB) Init() error {
	activateUuidExtentionQuery := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	b.Db.Exec(activateUuidExtentionQuery)
	createTypeQuery := fmt.Sprintf("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'game_result') THEN CREATE TYPE game_result AS ENUM ('%s','%s','%s'); END IF; END$$;", Draw, CityWin, MafiaWin)
	b.Db.Exec(createTypeQuery)
	err := b.Db.AutoMigrate(&Game{}, &PlayerGame{}, &User{})
	if err != nil {
		return err
	}
	return b.Db.SetupJoinTable(&Game{}, "Players", &PlayerGame{})
}
