package test

import (
	"fmt"
	db2 "github.com/GeorgeS1995/mafia_bot/internal/cfg/db"
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes Copy paste from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type DBTestSuite struct {
	suite.Suite
	db *db.MafiaDB
	Tx *gorm.DB
}

func (s *DBTestSuite) SetupSuite() {
	config, err := db2.NewMafiaBotDBConfig()
	if err != nil {
		log.Fatal("Can't parse config, ", err)
		return
	}

	dbObj, err := db.NewMafiaDB(*config)
	if err != nil {
		log.Fatal("Can't create db connection, ", err)
		return
	}

	testDBname := "test_" + config.DBname
	// https://stackoverflow.com/questions/54048774/how-to-create-a-postgres-database-using-gorm
	// check if db exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", testDBname)
	rs := dbObj.Db.Raw(stmt)
	if rs.Error != nil {
		log.Fatal("Can't check that test db exists, ", rs.Error)
		return
	}

	// if exists drop old db
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) > 0 {
		stmt = fmt.Sprintf("DROP DATABASE \"%s\";", testDBname)
		if rs = dbObj.Db.Exec(stmt); rs.Error != nil {
			log.Fatal("Can't drop old test db, ", rs.Error)
			return
		}

	}

	stmt = fmt.Sprintf("CREATE DATABASE \"%s\";", testDBname)
	rs = dbObj.Db.Exec(stmt)
	if rs.Error != nil {
		log.Fatal("Can't create test db, ", rs.Error)
		return
	}

	// close db connection
	sql, err := dbObj.Db.DB()
	_ = sql.Close()
	if err != nil {
		log.Fatal("Can't close connection to create test db, ", rs.Error)
		return
	}

	config.DBname = testDBname
	dbObj, err = db.NewMafiaDB(*config)
	if err != nil {
		log.Fatal("Can't create db connection, ", err)
		return
	}

	err = dbObj.Init()
	if err != nil {
		log.Fatal("Can't init db migration, ", err)
		return
	}

	s.db = dbObj
}

func (s *DBTestSuite) SetupTest() {
	s.Tx = s.db.Db.Begin()
}

func (s *DBTestSuite) TearDownSuite() {
	dbObj := s.db
	config, err := db2.NewMafiaBotDBConfig()
	if err != nil {
		log.Fatal("Can't parse config, ", err)
		return
	}

	testDBname := "test_" + config.DBname
	// https://stackoverflow.com/questions/54048774/how-to-create-a-postgres-database-using-gorm
	sqlDB, err := dbObj.Db.DB()
	if err != nil {
		log.Fatal("Can't, get sql db object from gorm db object", err)
		return
	}
	// Close
	err = sqlDB.Close()
	if err != nil {
		log.Fatal("Can't, get sql db object from gorm db object", err)
		return
	}

	dbObj, err = db.NewMafiaDB(*config)
	if err != nil {
		log.Fatal("Can't create db connection, ", err)
		return
	}
	stmt := fmt.Sprintf("DROP DATABASE \"%s\";", testDBname)
	rs := dbObj.Db.Exec(stmt)
	if rs.Error != nil {
		log.Fatal("Can't delete test db after test session, ", rs.Error)
		return
	}
}
