package db

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"os"
)

type MafiaBotDBConfig struct {
	Host     string
	User     string
	Password string
	DBname   string
	DSN      string
}

func NewMafiaBotDBConfig() (*MafiaBotDBConfig, error) {
	dsn_template := "host=%s user=%s password=%s dbname=%s port=5432"
	dbConfig := &MafiaBotDBConfig{}
	host := dbConfig.getHost()
	user, err := dbConfig.getUser()
	if err != nil {
		return dbConfig, err
	}
	password, err := dbConfig.getPassword()
	if err != nil {
		return dbConfig, err
	}
	dbname, err := dbConfig.getDBName()
	if err != nil {
		return dbConfig, err
	}
	dsn := fmt.Sprintf(dsn_template, host, user, password, dbname)
	dbConfig.Host = host
	dbConfig.User = user
	dbConfig.Password = password
	dbConfig.DBname = dbname
	dbConfig.DSN = dsn
	return dbConfig, nil
}

func (c *MafiaBotDBConfig) getHost() string {
	host := os.Getenv(fmt.Sprintf(common.ConfPrefix, "DB_HOST"))
	if host == "" {
		host = "mafia-db"
	}
	return host
}

// TODO get rid of copypast
func (c *MafiaBotDBConfig) getUser() (user string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "DB_USER")
	user = os.Getenv(env_name)
	if user == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return user, err
}

func (c *MafiaBotDBConfig) getPassword() (password string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "DB_PASSWORD")
	password = os.Getenv(env_name)
	if password == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return password, err
}

func (c *MafiaBotDBConfig) getDBName() (dbName string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "DB_NAME")
	dbName = os.Getenv(env_name)
	if dbName == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return dbName, err
}
