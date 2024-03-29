package db

import (
	"fmt"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/db"
	"github.com/GeorgeS1995/mafia_bot/test"
	"os"
	"testing"
)

func TestGetDSNSuccesfullEmptyHost(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_DB_USER")
		_ = os.Unsetenv("MAFIA_BOT_DB_PASSWORD")
		_ = os.Unsetenv("MAFIA_BOT_DB_NAME")
		_ = os.Unsetenv("MAFIA_BOT_DB_HOST")
	}()
	user := test.RandStringRunes(3)
	password := test.RandStringRunes(3)
	dbname := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DB_USER", user)
	_ = os.Setenv("MAFIA_BOT_DB_PASSWORD", password)
	_ = os.Setenv("MAFIA_BOT_DB_NAME", dbname)
	_ = os.Setenv("MAFIA_BOT_DB_HOST", "")

	dbConfig, err := db.NewMafiaBotDBConfig()

	expectedDSN := fmt.Sprintf("host=mafia-db user=%s password=%s dbname=%s port=5432", user, password, dbname)
	dsn := dbConfig.GetDSN()
	if dsn != expectedDSN {
		t.Fatalf("Dsn %v is not equal expected.", dsn)
	}
	if err != nil {
		t.Fatalf("Unexpected error while parsing db config: %v", err)
	}
}

func TestGetDSNSuccesfullWithHost(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_DB_HOST")
		_ = os.Unsetenv("MAFIA_BOT_DB_USER")
		_ = os.Unsetenv("MAFIA_BOT_DB_PASSWORD")
		_ = os.Unsetenv("MAFIA_BOT_DB_NAME")
	}()
	host := test.RandStringRunes(3)
	user := test.RandStringRunes(3)
	password := test.RandStringRunes(3)
	dbname := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DB_HOST", host)
	_ = os.Setenv("MAFIA_BOT_DB_USER", user)
	_ = os.Setenv("MAFIA_BOT_DB_PASSWORD", password)
	_ = os.Setenv("MAFIA_BOT_DB_NAME", dbname)

	dbConfig, err := db.NewMafiaBotDBConfig()

	expectedDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", host, user, password, dbname)
	dsn := dbConfig.GetDSN()
	if dsn != expectedDSN {
		t.Fatalf("Dsn %v is not equal expected.", dsn)
	}
	if err != nil {
		t.Fatalf("Unexpected error while parsing db config: %v", err)
	}
}

func TestGetDSNRequiredAttrError(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_DB_USER")
		_ = os.Unsetenv("MAFIA_BOT_DB_PASSWORD")
		_ = os.Unsetenv("MAFIA_BOT_DB_NAME")
	}()
	user := test.RandStringRunes(3)
	password := test.RandStringRunes(3)
	dbname := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_DB_USER", user)
	_ = os.Setenv("MAFIA_BOT_DB_PASSWORD", password)
	_ = os.Setenv("MAFIA_BOT_DB_NAME", dbname)

	for _, requiredAttr := range [3]string{"MAFIA_BOT_DB_USER", "MAFIA_BOT_DB_PASSWORD", "MAFIA_BOT_DB_NAME"} {
		t.Run(fmt.Sprintf("Current env %s", requiredAttr), func(t *testing.T) {
			_ = os.Unsetenv(requiredAttr)

			dbConfig, err := db.NewMafiaBotDBConfig()

			dsn := dbConfig.GetDSN()
			if dsn != "" {
				t.Fatalf("Dsn %v is not equal expected.", dsn)
			}
			var expectedErr common.MafiaBotParseMissingRequiredParamError
			expectedErr = common.MafiaBotParseMissingRequiredParamError{ParsedAttr: requiredAttr}
			expectedErrorMsg := expectedErr.Error()
			if err == nil || err.Error() != expectedErrorMsg {
				t.Fatalf("Unexpected error while parsing db config: %v", err)
			}
			_ = os.Setenv("MAFIA_BOT_DB_USER", user)
			_ = os.Setenv("MAFIA_BOT_DB_PASSWORD", password)
			_ = os.Setenv("MAFIA_BOT_DB_NAME", dbname)
		})
	}
}
