package component

import "github.com/bwmarrin/discordgo"

// Component is the base structure that must be aliased
// to create a new component.
//
// It holds basic metadata about the component
type Component struct {
    Name         string
    Description  string
    DmPermission bool
    Loaded       bool
    Enabled      bool
}

// LoadableComponent is the interface that allows a component to be registered.
// It must be implemented by new components to make it possible to register them
//
// The function returns the original component for use yb the registration
// system
type LoadableComponent interface {
    LoadComponent(discord *discordgo.Session) (*Component, error)
}
