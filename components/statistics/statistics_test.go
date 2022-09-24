package statistics

import (
	"github.com/lazybytez/jojo-discord-bot/test/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatisticsInit(t *testing.T) {
	result := helper.TestIfComponentIsRegistered(&C)

	assert.True(t, result)
}
