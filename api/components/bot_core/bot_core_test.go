package bot_core

import (
	"github.com/lazybytez/jojo-discord-bot/test/helper"
	"testing"
)

func TestBotCoreInit(t *testing.T) {
	helper.TestIfComponentIsRegistered(t, &C)
}
