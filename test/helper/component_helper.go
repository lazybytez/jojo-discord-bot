package helper

import (
	"github.com/lazybytez/jojo-discord-bot/api"
)

// TestIfComponentIsRegistered checks if the given component is registered.
// If the component is not registered. A test fail will be triggered.
func TestIfComponentIsRegistered(comp *api.Component) bool {
	for _, registeredComponent := range api.Components {
		if registeredComponent == comp {
			return true
		}
	}

	return false
}
