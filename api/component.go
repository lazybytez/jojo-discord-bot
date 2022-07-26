package api

import "github.com/bwmarrin/discordgo"

// Component is the base structure that must be aliased
// to create a new component.
//
// It holds basic metadata about the component
type Component struct {
    // Metadata
    Name         string
    Description  string
    DmPermission bool

    // State
    Loaded  bool
    Enabled bool

    // Lifecycle closures
    LoadComponent   func(discord *discordgo.Session) error
    UnloadComponent func(discord *discordgo.Session) error
}

// RegistrableComponent is the interface that allows a component to be registered.
type RegistrableComponent interface {
    RegisterComponent(discord *discordgo.Session) error
    UnregisterComponent(discord *discordgo.Session) error
}

// RegisterComponent is used by the component registration system that
// automatically calls the RegisterComponent method for all Component instances in
// the components.Components array.
func (c *Component) RegisterComponent(discord *discordgo.Session) error {
    err := c.LoadComponent(discord)

    if err != nil {
        return err
    }
    c.Loaded = true

    return nil
}

// UnregisterComponent is used by the component registration system that
// automatically calls the UnregisterComponent method for all Component instances in
// the components.Components array.
//
// It is used to give components the ability to gracefully shutdown.
func (c *Component) UnregisterComponent(discord *discordgo.Session) error {
    err := c.UnloadComponent(discord)

    if err != nil {
        return err
    }
    c.Loaded = false

    return nil
}
