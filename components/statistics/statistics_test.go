package statistics

import (
	"github.com/lazybytez/jojo-discord-bot/test/helper"
	"testing"
)

func TestStatisticsInit(t *testing.T) {
	helper.TestIfComponentIsRegistered(t, &C)
}
