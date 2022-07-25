package components

import (
    "github.com/bwmarrin/discordgo"
    "github.com/lazybytez/jojo-discord-bot/components/pingpong"
    "github.com/lazybytez/jojo-discord-bot/internal/component"
    "github.com/rs/zerolog/log"
)

// Components contains all components that should be available.
//
// Enabled components should be registered here.
var Components = []component.LoadableComponent{
    pingpong.ComponentInstance,
}

// RegisterComponents handles the initialization of
// all components listed in the Components array.
//
// When it is not possible to register a component,
// an error will be printed into the log.
// The application will continue to run as nothing happened.
func RegisterComponents(discord *discordgo.Session) {
    for _, comp := range Components {
        comp, err := comp.LoadComponent(discord)
        if nil != err {
            notifyComponentLoadFailed(comp, err)
            continue
        }
        comp.Loaded = true
    }
}

// notifyComponentLoadFailed prints a message to the log
// that contains information of a failure when loading a component.
//
// This function is used by RegisterComponents to notify component
// loading failures.
func notifyComponentLoadFailed(comp *component.Component, err error) {
    log.Warn().Msgf("Failed to load component with name \"%v\": %v", comp.Name, err.Error())
}
