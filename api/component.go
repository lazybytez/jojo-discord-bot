package api

import (
	"github.com/bwmarrin/discordgo"
)

// LifecycleHooks allow to specify functions that should be called
// when components get loaded and unloaded.
//
// The defined functions allow some way of initializing a owner.
type LifecycleHooks struct {
	LoadComponent   func(discord *discordgo.Session) error
	UnloadComponent func(discord *discordgo.Session) error
}

// State holds the state of a owner.
// This includes states like:
//   - is the owner enabled?
//   - is the owner currently loaded?
type State struct {
	Loaded  bool
	Enabled bool
}

// Component is the base structure that must be aliased
// to create a new owner.
//
// It holds basic metadata about the owner
type Component struct {
	// Metadata
	Name         string
	Description  string
	DmPermission bool

	// State
	State State

	// Lifecycle hooks
	Lifecycle LifecycleHooks

	// Utilities
	// These are private and only managed by the API system.
	// Their initialization happens through call to the methods
	// used to get them (Example: logger -> Component.Logger()).
	logger         Logger
	handlerManager ComponentHandlerManager
	discord        *discordgo.Session
}

// RegistrableComponent is the interface that allows a owner to be
// initialized and registered.
type RegistrableComponent interface {
	RegisterComponent(discord *discordgo.Session) error
	UnregisterComponent(discord *discordgo.Session) error
}

// RegisterComponent is used by the owner registration system that
// automatically calls the RegisterComponent method for all Component instances in
// the components.Components array.
func (c *Component) RegisterComponent(discord *discordgo.Session) error {
	c.discord = discord

	err := c.Lifecycle.LoadComponent(discord)

	if err != nil {
		return err
	}
	c.State.Loaded = true

	return nil
}

// UnregisterComponent is used by the owner registration system that
// automatically calls the UnregisterComponent method for all Component instances in
// the components.Components array.
//
// It is used to give components the ability to gracefully shutdown.
func (c *Component) UnregisterComponent(discord *discordgo.Session) error {
	err := c.Lifecycle.UnloadComponent(discord)

	if err != nil {
		return err
	}
	c.State.Loaded = false

	return nil
}
