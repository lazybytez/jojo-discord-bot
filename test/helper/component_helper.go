package helper

import (
	"github.com/lazybytez/jojo-discord-bot/api"
	"testing"
)

// TestIfComponentIsRegistered checks if the given component is registered.
// If the component is not registered. A test fail will be triggered.
func TestIfComponentIsRegistered(t *testing.T, comp *api.Component) {
	for _, registeredComponent := range api.Components {
		if registeredComponent == comp {
			return
		}
	}

	t.Fail()
}
