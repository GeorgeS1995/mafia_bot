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
}

func NewMafiaBotDBConfig() (*MafiaBotDBConfig, error) {
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
	dbConfig.Host = host
	dbConfig.User = user
	dbConfig.Password = password
	dbConfig.DBname = dbname
	return dbConfig, nil
}

func (c *MafiaBotDBConfig) GetDSN() string {
	if c.Host == "" || c.User == "" || c.Password == "" || c.DBname == "" {
		return ""
	}
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=5432"
	return fmt.Sprintf(dsnTemplate, c.Host, c.User, c.Password, c.DBname)
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
	envName := fmt.Sprintf(common.ConfPrefix, "DB_USER")
	user = os.Getenv(envName)
	if user == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return user, err
}

func (c *MafiaBotDBConfig) getPassword() (password string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "DB_PASSWORD")
	password = os.Getenv(envName)
	if password == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return password, err
}

func (c *MafiaBotDBConfig) getDBName() (dbName string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "DB_NAME")
	dbName = os.Getenv(envName)
	if dbName == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return dbName, err
}
