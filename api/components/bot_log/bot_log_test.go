package bot_log

import (
	"github.com/lazybytez/jojo-discord-bot/test/helper"
	"testing"
)

func TestBotLogInit(t *testing.T) {
	helper.TestIfComponentIsRegistered(t, &C)
}
