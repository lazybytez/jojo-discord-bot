/*
 * JOJO Discord Bot - An advanced multi-purpose discord bot
 * Copyright (C) 2022 Lazy Bytez (Elias Knodel, Pascal Zarrad)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/util"
	"reflect"
	"sync"
)

// Management tools that allow registering and unregistering
// DiscordGo event handlers.
//
// This acts as a wrapper structure to allow the bot framework
// to take some actions in-between like logging handler
// registration or injecting custom logic like the component status checks.
//

// componentHandlerMap is a wrapper that holds the map that contains
// the Handler name -> Handler mapping.
type componentHandlerMap struct {
	sync.RWMutex
	handlers map[string]*AssignedEventHandler
}

// handlerComponentMapping is a map that holds references to function's
// that are registered as handlers and maps them to their owning components.
//
// Note that the key of the map equals the handler name.
// The handler name is always the component name in snake-case
// followed by the Handler name defined by the developer.
//
// The reason for doing this is, to allow future adjustments of handlers,
// by holding a reference that can be edited.
var handlerComponentMapping = componentHandlerMap{
	handlers: make(map[string]*AssignedEventHandler),
}

// AssignedEventHandler holds the necessary data to handle events
// in the bot.
// The originalHandler should only be set during creation of the handler.
// The handler can be swapped with some new function that wraps the original function.
// This allows to add custom logic like the component status checks.
//
// Important: The handler can only being wrapped before it is registered
// in DiscordGo. Wrappers are meant for system features of the bot like
// component status checks or permissions. Those features are injected during the
// call of the ComponentHandlerContainer.Register and ComponentHandlerContainer.RegisterOnce methods.
type AssignedEventHandler struct {
	name            string
	component       *Component
	originalHandler interface{}
	handler         interface{}

	unregister func()
}

// AssignedEventHandlerAccess is an interface that provides access to
// some fields of AssignedEventHandler through getters.
//
// This ensures those sensitive values don't get overridden.
type AssignedEventHandlerAccess interface {
	GetName() string
	GetComponent() *Component
	GetHandler() interface{}
}

// GetName returns the name of the Handler assigned to the AssignedEventHandler
func (ase *AssignedEventHandler) GetName() string {
	return ase.name
}

// GetComponent returns the component that owns the Handler assigned to the AssignedEventHandler
func (ase *AssignedEventHandler) GetComponent() *Component {
	return ase.component
}

// GetHandler returns the original Handler function assigned to the AssignedEventHandler
func (ase *AssignedEventHandler) GetHandler() interface{} {
	return ase.handler
}

// ComponentHandlerContainer is a helper that eases registration of event handlers.
// The methods provided by ComponentHandlerContainer include
// automated logging and wrappers to apply some standards to all handlers.
type ComponentHandlerContainer struct {
	owner *Component
}

// ComponentHandlerManager is the interface that defines all
// methods for event handler management.
//
// The functions of the interface allow registration and removal
// of Discord event handlers.
//
// Under the hood, DiscordGo event handlers are used.
type ComponentHandlerManager interface {
	// Register can be used to register a new AssignedEventHandler.
	//
	// The passed handler function will be:
	//   1. registered in DiscordGo as a Handler
	//   2. wrapped with some internal logic of the bot
	//   3. saved with a name that allows to retrieve it later on
	//
	// The passed name for the handler is concatenated with the name of the
	// component that owns the handler (separated by underscore).
	//
	// The handler must have the same format as when a handler is registered in
	// plain DiscordGo. See the documentation about discordgo.AddHandler
	// for additional information.
	//
	// In general, the common format for a handler function is:
	//   func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
	Register(name string, handler interface{}) (string, bool)

	// RegisterOnce registers an event handler as a one-time event Handler.
	// The registered handler will be removed after being triggered once.
	//
	//
	// The passed handler function will be:
	//   1. registered in DiscordGo as a Handler
	//   2. wrapped with some internal logic of the bot
	//   3. saved with a name that allows to retrieve it later on
	//
	// The passed name for the handler is concatenated with the name of the
	// component that owns the handler (separated by underscore).
	//
	// The handler must have the same format as when a Handler is registered in
	// plain DiscordGo. See the documentation about discordgo.AddHandler
	// for additional information.
	//
	// In general, the common format for a handler function is:
	//   func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
	RegisterOnce(name string, handler interface{}) (string, bool)

	// Unregister removes the handler with the given name (if existing) from
	// the registered handlers.
	//
	// If the specified handler does not exist, an error will be returned.
	Unregister(name string) error

	UnregisterAll()
}

// HandlerManager returns the management interface for event handlers.
//
// It allows the general management of event handlers.
//
// It should be always used when event handlers to listen to
// Discord events are necessary. It natively handles stuff like logging,
// component status and permissions.
func (c *Component) HandlerManager() ComponentHandlerManager {
	if nil == c.handlerManager {
		c.handlerManager = &ComponentHandlerContainer{
			owner: c,
		}
	}

	return c.handlerManager
}

// === Handler registration

// Register can be used to register a new AssignedEventHandler.
//
// The passed handler function will be:
//  1. registered in DiscordGo as a Handler
//  2. wrapped with some internal logic of the bot
//  3. saved with a name that allows to retrieve it later on
//
// The passed name for the handler is concatenated with the name of the
// component that owns the handler (separated by underscore).
//
// The handler must have the same format as when a handler is registered in
// plain DiscordGo. See the documentation about discordgo.AddHandler
// for additional information.
//
// In general, the common format for a handler function is:
//
//	func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
func (c *ComponentHandlerContainer) Register(name string, handler interface{}) (string, bool) {
	handlerName := GetHandlerName(c.owner, name)

	if _, ok := GetHandler(handlerName); ok {
		c.owner.Logger().Err(fmt.Errorf(
			"an handler for component \"%v\" with name \"%v\" is already registered (ID: \"%v\")",
			c.owner.Name,
			name,
			handlerName),
			"Failed to register handler using handler management!")

		return handlerName, false
	}

	assignedEvent := &AssignedEventHandler{
		name:            handlerName,
		component:       c.owner,
		originalHandler: handler,
		handler:         handler,
	}
	c.addComponentHandler(handlerName, assignedEvent)

	// Apply system wrappers that add some logic to the original handler
	c.wrapWithComponentStatusHandler(handlerName)

	c.addDiscordGoHandler(assignedEvent)

	c.owner.Logger().Info("Handler with name \"%v\" for component \"%v\" has been registered! (ID: \"%v\")",
		name,
		c.owner.Name,
		handlerName)

	return handlerName, true
}

// addDiscordGoHandler registers the handler in the standard DiscordGo session.
// The handler function is before wrapped with some function that ensures the proper
// event is used by DiscordGo, as the handler wrappers are not strictly typed.
func (c *ComponentHandlerContainer) addDiscordGoHandler(assignedEvent *AssignedEventHandler) {
	originalType := reflect.TypeOf(assignedEvent.originalHandler)
	typedHandler := reflect.MakeFunc(originalType, func(args []reflect.Value) (results []reflect.Value) {
		return reflect.ValueOf(assignedEvent.handler).Call(args)
	})

	c.owner.discord.AddHandler(typedHandler.Interface())
}

// === One-Time Handlers

// RegisterOnce registers an event handler as a one-time event Handler.
// The registered handler will be removed after being triggered once.
//
// The passed handler function will be:
//  1. registered in DiscordGo as a Handler
//  2. wrapped with some internal logic of the bot
//  3. saved with a name that allows to retrieve it later on
//
// The passed name for the handler is concatenated with the name of the
// component that owns the handler (separated by underscore).
//
// The handler must have the same format as when a Handler is registered in
// plain DiscordGo. See the documentation about discordgo.AddHandler
// for additional information.
//
// In general, the common format for a handler function is:
//
//	func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
func (c *ComponentHandlerContainer) RegisterOnce(
	name string,
	handler interface{},
) (string, bool) {
	handlerName := GetHandlerName(c.owner, name)

	if _, ok := GetHandler(handlerName); ok {
		c.owner.Logger().Err(fmt.Errorf(
			"an Handler for component \"%v\" with name \"%v\" is already registered (ID: \"%v\")",
			c.owner.Name,
			name,
			handlerName),
			"Failed to register one-time handler!")

		return handlerName, false
	}

	assignedEvent := &AssignedEventHandler{
		name:            handlerName,
		component:       c.owner,
		originalHandler: handler,
		handler:         handler,
	}
	c.addComponentHandler(handlerName, assignedEvent)

	// Apply system wrappers that add some logic to the original handler
	c.wrapWithComponentStatusHandler(handlerName)
	c.wrapOneTimeHandler(handlerName)

	c.addDiscordGoOnceTimeHandler(assignedEvent)

	c.owner.Logger().Info(
		"One-time Handler with name \"%v\" for component \"%v\" has been registered! (ID: \"%v\")",
		name,
		c.owner.Name,
		handlerName)

	return handlerName, true
}

// addDiscordGoOnceTimeHandler registers the handler in the standard DiscordGo session.
// The handler function is before wrapped with some function that ensures the proper
// event is used by DiscordGo, as the handler wrappers are not strictly typed.
// The handler is registered as a once time handler.
func (c *ComponentHandlerContainer) addDiscordGoOnceTimeHandler(assignedEvent *AssignedEventHandler) {
	originalType := reflect.TypeOf(assignedEvent.originalHandler)
	typedHandler := reflect.MakeFunc(originalType, func(args []reflect.Value) (results []reflect.Value) {
		return reflect.ValueOf(assignedEvent.handler).Call(args)
	})

	c.owner.discord.AddHandlerOnce(typedHandler.Interface())
}

// === Handler removal

// Unregister removes the Handler with the given name (if existing) from
// the registered handlers.
//
// If the specified Handler does not exist, an error will be returned.
func (c *ComponentHandlerContainer) Unregister(name string) error {
	handlerName := GetHandlerName(c.owner, name)
	handler, ok := GetHandler(handlerName)

	if !ok {
		return fmt.Errorf(
			"there is no handler called \"%v\" registered that could be unregistered",
			handlerName)
	}

	handler.unregister()
	removeComponentHandler(handlerName)

	return nil
}

// UnregisterAll takes care of unregistering all handlers
// attached to the component that owns the ComponentHandlerContainer
func (c *ComponentHandlerContainer) UnregisterAll() {
	for _, handler := range handlerComponentMapping.handlers {
		if handler.component == c.owner {
			_ = c.Unregister(handler.name)
		}
	}
}

// === Handler management

// addComponentHandler adds a new Handler to the registered handlers.
//
// Note that adding a Handler with a name that is already in the map
// will override the existing Handler (but not unregister it from DiscordGo!).
func (c *ComponentHandlerContainer) addComponentHandler(name string, handler *AssignedEventHandler) {
	handlerComponentMapping.Lock()
	defer handlerComponentMapping.Unlock()

	handlerComponentMapping.handlers[name] = handler
}

// removeComponentHandler removes the Handler with the specified name
// from the registered handlers.
//
// It also removes all decorators that are assigned to the
// Handler.
func removeComponentHandler(name string) {
	handlerComponentMapping.Lock()
	defer handlerComponentMapping.Unlock()

	delete(handlerComponentMapping.handlers, name)
}

// GetHandler returns a Handler by its fully qualified name (id).
// The required ID can be obtained using GetHandlerName.
func GetHandler(name string) (*AssignedEventHandler, bool) {
	handlerComponentMapping.RLock()
	defer handlerComponentMapping.RUnlock()

	handler, ok := handlerComponentMapping.handlers[name]

	return handler, ok
}

// GetHandlerName returns the name of a Handler for a component.
//
// It acts as the auto-formatter that should be used to retrieve
// Handler names.
func GetHandlerName(c *Component, name string) string {
	return util.StringToSnakeCase(fmt.Sprintf("%v_%v", c.Name, name))
}

// === Handler wrappers

// wrapWithComponentStatusHandler wraps the original handler function and
// ensures that the original handler function is only called when
// the owning component is enabled.
func (c *ComponentHandlerContainer) wrapWithComponentStatusHandler(name string) {
	assignedEventHandler, ok := GetHandler(name)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed wrap original handler with name \"%v\" of component \"%v\" to prevent it"+
				" being executed when component is disabled",
			name,
			c.owner.Name),
			"Failed to apply component status handler!")
	}

	originalHandler := assignedEventHandler.handler
	assignedEventHandler.handler = func(
		session *discordgo.Session,
		event interface{},
	) {
		comp := assignedEventHandler.GetComponent()
		guildId := getGuildIdFromEventInterface(event)

		if IsComponentEnabled(comp, guildId) {
			reflect.ValueOf(originalHandler).Call([]reflect.Value{
				reflect.ValueOf(session),
				reflect.ValueOf(event),
			})
		}
	}
}

// wrapOneTimeHandler is a wrapper appended functions used in one-time
// event handlers. It ensures that the executed Handler is removed properly.
func (c *ComponentHandlerContainer) wrapOneTimeHandler(name string) {
	assignedEventHandler, ok := GetHandler(name)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed wrap original handler with name \"%v\" of component \"%v\" to remove it when executed once",
			name,
			c.owner.Name),
			"Failed to apply component one-time removal handler!")
	}

	originalHandler := assignedEventHandler.handler
	assignedEventHandler.handler = func(
		session *discordgo.Session,
		event interface{},
	) {
		removeComponentHandler(assignedEventHandler.name)

		reflect.ValueOf(originalHandler).Call([]reflect.Value{
			reflect.ValueOf(session),
			reflect.ValueOf(event),
		})
	}
}
