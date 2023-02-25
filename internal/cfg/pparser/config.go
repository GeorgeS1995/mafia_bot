package pparser

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"os"
	"strconv"
)

type MafiaBotPparserConfig struct {
	CSRF                  string
	CSRFCookie            string
	PolemicaHost          string
	Login                 string
	Password              string
	ParseHistoryTaskDelay int
}

func NewMafiaBotPparserConfig() (*MafiaBotPparserConfig, error) {
	pparserConfig := &MafiaBotPparserConfig{}
	csfr, err := pparserConfig.GetCSRF()
	if err != nil {
		return pparserConfig, err
	}
	csrfCookie, err := pparserConfig.GetCSRFCookie()
	if err != nil {
		return pparserConfig, err
	}
	polemicaHost, err := pparserConfig.GetPolemicaHost()
	if err != nil {
		return pparserConfig, err
	}
	polemicaLogin, err := pparserConfig.GetLogin()
	if err != nil {
		return pparserConfig, err
	}
	polemicaPassword, err := pparserConfig.GetPassword()
	if err != nil {
		return pparserConfig, err
	}
	polemicaParserDelay, err := pparserConfig.GetParseHistoryTaskDelay()
	if err != nil {
		return pparserConfig, err
	}
	pparserConfig.CSRF = csfr
	pparserConfig.CSRFCookie = csrfCookie
	pparserConfig.PolemicaHost = polemicaHost
	pparserConfig.Login = polemicaLogin
	pparserConfig.Password = polemicaPassword
	pparserConfig.ParseHistoryTaskDelay = polemicaParserDelay
	return pparserConfig, nil
}

func (c *MafiaBotPparserConfig) GetCSRF() (csrf string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "CSRF")
	csrf = os.Getenv(envName)
	if csrf == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return csrf, err
}

func (c *MafiaBotPparserConfig) GetCSRFCookie() (csrfCookie string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "CSRF_COOKIE")
	csrfCookie = os.Getenv(envName)
	if csrfCookie == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return csrfCookie, err
}

func (c *MafiaBotPparserConfig) GetPolemicaHost() (polemicaHost string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "POLEMICA_HOST")
	polemicaHost = os.Getenv(envName)
	if polemicaHost == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return polemicaHost, err
}

func (c *MafiaBotPparserConfig) GetLogin() (login string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "POLEMICA_LOGIN")
	login = os.Getenv(envName)
	if login == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return login, err
}

func (c *MafiaBotPparserConfig) GetPassword() (password string, err error) {
	envName := fmt.Sprintf(common.ConfPrefix, "POLEMICA_PASSWORD")
	password = os.Getenv(envName)
	if password == "" {
		err = &common.MafiaBotParseMissingRequiredParamError{ParsedAttr: envName}
	}
	return password, err
}

func (c *MafiaBotPparserConfig) GetParseHistoryTaskDelay() (int, error) {
	envName := fmt.Sprintf(common.ConfPrefix, "POLEMICA_PARSE_HISTORY_TASK_DELAY")
	delay := 1000 * 60 * 60 // default one check per hour
	delayStr := os.Getenv(envName)
	var err error
	if delayStr != "" {
		delay, err = strconv.Atoi(delayStr)
		if err != nil {
			err = &common.MafiaBotParseTypeError{ParsedAttr: envName}

		}
	}
	return delay, err
}
