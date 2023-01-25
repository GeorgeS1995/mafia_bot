package pparser

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"os"
)

type MafiaBotPparserConfig struct {
	CSRF         string
	CSRFCookie   string
	PolemicaHost string
	Login        string
	Password     string
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
	pparserConfig.CSRF = csfr
	pparserConfig.CSRFCookie = csrfCookie
	pparserConfig.PolemicaHost = polemicaHost
	pparserConfig.Login = polemicaLogin
	pparserConfig.Password = polemicaPassword
	return pparserConfig, nil
}

func (c *MafiaBotPparserConfig) GetCSRF() (csrf string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "CSRF")
	csrf = os.Getenv(env_name)
	if csrf == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return csrf, err
}

func (c *MafiaBotPparserConfig) GetCSRFCookie() (csrfCookie string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "CSRF_COOKIE")
	csrfCookie = os.Getenv(env_name)
	if csrfCookie == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return csrfCookie, err
}

func (c *MafiaBotPparserConfig) GetPolemicaHost() (polemicaHost string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "POLEMICA_HOST")
	polemicaHost = os.Getenv(env_name)
	if polemicaHost == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return polemicaHost, err
}

func (c *MafiaBotPparserConfig) GetLogin() (login string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "POLEMICA_LOGIN")
	login = os.Getenv(env_name)
	if login == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return login, err
}

func (c *MafiaBotPparserConfig) GetPassword() (password string, err error) {
	env_name := fmt.Sprintf(common.ConfPrefix, "POLEMICA_PASSWORD")
	password = os.Getenv(env_name)
	if password == "" {
		err = &common.MafiaBotParseError{ParsedAttr: env_name}
	}
	return password, err
}
