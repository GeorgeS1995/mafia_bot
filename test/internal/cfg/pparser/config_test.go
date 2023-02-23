package pparser

import (
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/common"
	"github.com/GeorgeS1995/mafia_bot/internal/cfg/pparser"
	"github.com/GeorgeS1995/mafia_bot/test"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestGetParseHistoryTaskDelayTypeError(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_POLEMICA_PARSE_HISTORY_TASK_DELAY")
	}()
	wrongDelay := test.RandStringRunes(3)
	_ = os.Setenv("MAFIA_BOT_POLEMICA_PARSE_HISTORY_TASK_DELAY", wrongDelay)
	config := &pparser.MafiaBotPparserConfig{}

	_, err := config.GetParseHistoryTaskDelay()

	if _, ok := err.(*common.MafiaBotParseTypeError); !ok {
		t.Fatalf("Wrong error type: %s", err)
	}
}

func TestGetParseHistoryTaskDelayDefaultOK(t *testing.T) {
	config := &pparser.MafiaBotPparserConfig{}

	delay, err := config.GetParseHistoryTaskDelay()

	if err != nil || delay != 1000*60*60 {
		t.Fatalf("Can't set default value.\n Error: %s\n, Value: %d\n", err, delay)
	}
}

func TestGetParseHistoryTaskDelayOK(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("MAFIA_BOT_POLEMICA_PARSE_HISTORY_TASK_DELAY")
	}()
	rand.Seed(time.Now().UnixNano())
	envDelay := rand.Intn(100001)
	_ = os.Setenv("MAFIA_BOT_POLEMICA_PARSE_HISTORY_TASK_DELAY", strconv.Itoa(envDelay))

	config := &pparser.MafiaBotPparserConfig{}

	delay, err := config.GetParseHistoryTaskDelay()

	if err != nil || delay != envDelay {
		t.Fatalf("Can't set default value.\n Error: %s\n, Value: %d\n, Expected delay: %d\n", err, delay, envDelay)
	}
}
