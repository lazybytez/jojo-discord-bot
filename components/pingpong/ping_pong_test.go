package pingpong

import (
	"github.com/lazybytez/jojo-discord-bot/test/helper"
	"testing"
)

func TestPingPongInit(t *testing.T) {
	helper.TestIfComponentIsRegistered(t, &C)
}
