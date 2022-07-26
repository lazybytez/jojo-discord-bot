package components

import (
    "github.com/bwmarrin/discordgo"
    "github.com/lazybytez/jojo-discord-bot/api"
    "github.com/lazybytez/jojo-discord-bot/api/log"
    "github.com/lazybytez/jojo-discord-bot/components/pingpong"
)

// LogComponentRegistry is the custom component name used
// to identify log messages from the component management system
const LogComponentRegistry = "Component Manager"

// Components contains all components that should be available.
//
// Enabled components should be registered here.
var Components = []*api.Component{
    pingpong.ComponentInstance,
}

// RegisterComponents handles the initialization of
// all components listed in the Components array.
//
// When it is not possible to register a component,
// an error will be printed into the log.
// The application will continue to run as nothing happened.
func RegisterComponents(discord *discordgo.Session) {
    log.LogInfo(LogComponentRegistry, "Starting component load sequence...")
    for _, comp := range Components {
        if nil == comp.LoadComponent {
            log.LogDebug(LogComponentRegistry, "Component \"%v\" does not have an load callback, not loading it!", comp.Name)
            continue
        }

        if !comp.Enabled {
            log.LogInfo(LogComponentRegistry, "Component \"%v\" is not enabled, skipping!", comp.Name)
            continue
        }

        log.LogInfo(LogComponentRegistry, "Loading component \"%v\"...", comp.Name)
        err := comp.RegisterComponent(discord)
        if nil != err {
            log.LogWarn(LogComponentRegistry, "Failed to load component with name \"%v\": %v", comp.Name, err.Error())
            continue
        }
        log.LogInfo(LogComponentRegistry, "Successfully loaded component \"%v\"!", comp.Name)
    }
    log.LogInfo(LogComponentRegistry, "Component load sequence completed!")
}

// UnloadComponents iterates through all registered api.Component
// registered in the Components array and calls their UnloadComponent
// function.
//
// If an api.Component does not have an UnloadComponent function defined,
// it will be ignored.
func UnloadComponents(discord *discordgo.Session) {
    log.LogInfo(LogComponentRegistry, "Starting component unload sequence...")
    for _, comp := range Components {
        if nil == comp.UnloadComponent {
            log.LogDebug(LogComponentRegistry, "Component \"%v\" does not have an unload callback, skipping!", comp.Name)
            continue
        }

        if !comp.Loaded {
            log.LogWarn(LogComponentRegistry, "Component \"%v\" has not been loaded, skipping!", comp.Name)
            continue
        }

        log.LogInfo(LogComponentRegistry, "Unloading component \"%v\"...", comp.Name)
        err := comp.UnregisterComponent(discord)
        if nil != err {
            log.LogWarn(LogComponentRegistry, "Failed to unload component with name \"%v\": %v", comp.Name, err.Error())
            continue
        }
        log.LogInfo(LogComponentRegistry, "Successfully unloaded component \"%v\"!", comp.Name)
    }
    log.LogInfo(LogComponentRegistry, "Unload sequence completed!")
}
