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
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/util"
	"reflect"
	"sync"
)

// Management tools that allow registering, unregistering
// and decorating DiscordGo event handlers.
//
// This acts as a wrapper structure to allow the bot framework
// to take some actions in-between like logging handler
// registration.
//
// Also, this allows to decorate existing handlers and add
// custom functionality to them by intercepting their call.
// This is form example useful for the handlers that are removed after
// being called once.
// With the decorator feature, we can clean up out own handler list
// when the event handler is called.
//
// Note: This might not be the best way of implementing those features.
//       The implementation should be checked and reworked later on when
//       we know better ways of implementing this stuff.
//       The current interface should allow keeping the API,
//       even when the implementation details change.

// === Component structs, interfaces and internal data management

// componentHandlerMap is a wrapper that holds the map that contains
// the handler name -> handler mapping. It embeds sync.RWMutex
// to allow support for concurrency.
type componentHandlerMap struct {
	sync.RWMutex
	handlers map[string]*AssignedEventHandler
}

// handlerComponentMapping is a map that holds references to function's
// that are registered as handlers and maps them to their owning components.
//
// Note that the key of the map equals the handler name.
// The handler name is always the owner name in snake-case
// followed by the handler name defined by the developer.
//
// The reason for doing this is, to allow future adjustments of handlers,
// by holding a reference that can be edited.
var handlerComponentMapping = componentHandlerMap{
	handlers: make(map[string]*AssignedEventHandler),
}

// AssignedEventHandler holds the necessary data to handle events
// in the bot. It is used to allow features like decorating handlers
// by caching the original handler function.
// The original handler function is then replaced with a proxy function
// that allows us to add special features.
type AssignedEventHandler struct {
	name       string
	component  *Component
	handler    interface{}
	decorators *decoratorChain

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

// GetName returns the name of the handler assigned to the AssignedEventHandler
func (ase *AssignedEventHandler) GetName() string {
	return ase.name
}

// GetComponent returns the component that owns the handler assigned to the AssignedEventHandler
func (ase *AssignedEventHandler) GetComponent() *Component {
	return ase.component
}

// GetHandler returns the original handler function assigned to the AssignedEventHandler
func (ase *AssignedEventHandler) GetHandler() interface{} {
	return ase.handler
}

// ComponentHandlerContainer is a helper that eases registration of event handlers
// while establishing a standard for that task that includes
// automated logging and adds the ability to allow some weird
// stuff in the future.
type ComponentHandlerContainer struct {
	owner *Component
}

// ComponentHandlerManager is the interface that defines all
// methods for event management.
//
// The functions of the interface allow registration, decoration
// and removal of Discord event handlers.
//
// Under the hood, DiscordGo event handlers are used.
type ComponentHandlerManager interface {
	// Register can be used to register a new discordgo.AssignedEventHandler.
	//
	// The passed handler function will be:
	//   1. registered in DiscordGo as a handler
	//   2. prepared to allow decorations
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
	Register(name string, handler interface{}) (string, error)

	// RegisterOnce registers an event handler as a one-time event handler.
	// The registered handler will be removed after being triggered once.
	//
	//
	// The passed handler function will be:
	//   1. registered in DiscordGo as a handler
	//   2. prepared to allow decorations
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
	RegisterOnce(name string, handler interface{}) (string, error)

	// Unregister removes the handler with the given name (if existing) from
	// the registered handlers.
	//
	// If the specified handler does not exist, an error will be returned.
	Unregister(name string) error

	// AddDecorator allows to register a new decorator for an event
	// handler.
	//
	// Decorators will be called one after another until a decorator does not call
	// a following one (this causes the event handling to be cancelled).
	//
	// Decorator functions should have the following format:
	//   func (
	//      assignedEvent AssignedEventHandler,
	//      originalHandler interface{},
	//      session *discordgo.Session,
	//      event interface{}
	//   )
	//
	// Unclear params:
	//   - event must be the event that is handled by the original handler
	//   - originalHandler must be the type of the original handlers function
	AddDecorator(name string, decorator interface{}) error

	addDiscordGoHandler(assignedEvent *AssignedEventHandler)
	addComponentHandler(name string, handler *AssignedEventHandler)
}

// HandlerManager returns the management interface for event handlers.
//
// It allows the registration, decoration and general
// management of event handlers.
//
// It should be always used when event handlers to listen  for
// Discord events are necessary. It natively handles stuff like logging
// event handler status.
func (c *Component) HandlerManager() ComponentHandlerManager {
	if nil == c.handlerManager {
		c.handlerManager = &ComponentHandlerContainer{
			owner: c,
		}
	}

	return c.handlerManager
}

// === Handler registration

// Register can be used to register a new discordgo.AssignedEventHandler.
//
// The passed handler function will be:
//   1. registered in DiscordGo as a handler
//   2. prepared to allow decorations
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
func (c *ComponentHandlerContainer) Register(name string, handler interface{}) (string, error) {
	handlerName := GetHandlerName(c.owner, name)

	if _, ok := GetHandler(handlerName); ok {
		return handlerName, fmt.Errorf(
			"an handler for component \"%v\" with name \"%v\" is already registered (ID: \"%v\")",
			c.owner.Name,
			name,
			handlerName)
	}

	assignedEvent := &AssignedEventHandler{
		name:      handlerName,
		component: c.owner,
		handler:   handler,
	}
	c.addComponentHandler(handlerName, assignedEvent)

	c.addDiscordGoHandler(assignedEvent)

	c.owner.Logger().Info("Handler with name \"%v\" for component \"%v\" has been registered! (ID: \"%v\")",
		name,
		c.owner.Name,
		handlerName)

	return handlerName, nil
}

// createHandlerProxy creates a closure that acts as a proxy for the original handler function.
// This allows us to delegate the triggered event to the callOriginalHandler function.
// The callOriginalHandler than allows to intercept the original event handler with
// custom decorator functions. This can be useful in some edge-cases.
func createHandlerProxy(handler *AssignedEventHandler) func(args []reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		decorators := handler.decorators

		if nil == decorators || nil == decorators.Obtain() {
			decorateArguments := append([]reflect.Value{reflect.ValueOf(handler)}, args...)
			reflect.ValueOf(callOriginalHandler).Call(decorateArguments)

			return []reflect.Value{}
		}

		callOriginal := false
		currentDecorator := decorators.Obtain()
		for currentDecorator != nil {
			decoratorHandlerWrapper := func(args []reflect.Value) []reflect.Value {
				if nil == currentDecorator || nil == currentDecorator.next {
					callOriginal = true
					currentDecorator = nil

					return []reflect.Value{}
				}

				return []reflect.Value{}
			}

			decoratorHandlerEventHandler := reflect.MakeFunc(reflect.TypeOf(handler.handler), decoratorHandlerWrapper)

			loopHandler := currentDecorator

			decorateArguments := append([]reflect.Value{reflect.ValueOf(handler)}, args...)
			decorateArguments = append(decorateArguments, decoratorHandlerEventHandler)

			currentDecorator = currentDecorator.next
			reflect.ValueOf(loopHandler.value).Call(decorateArguments)
		}

		if !callOriginal {
			return []reflect.Value{}
		}

		decorateArguments := append([]reflect.Value{reflect.ValueOf(handler)}, args...)
		reflect.ValueOf(callOriginalHandler).Call(decorateArguments)

		return []reflect.Value{}
	}
}

// callOriginalHandler is the internal replacement used when an event happens in discord.
func callOriginalHandler(assignedEvent *AssignedEventHandler, session *discordgo.Session, e interface{}) {
	sessionRef := reflect.ValueOf(session)
	eRef := reflect.ValueOf(e)

	reflect.ValueOf(assignedEvent.handler).Call([]reflect.Value{sessionRef, eRef})
}

// addDiscordGoHandler generates a handler proxy and registers it for DiscordGo.
func (c *ComponentHandlerContainer) addDiscordGoHandler(assignedEvent *AssignedEventHandler) {
	handlerProxy := createHandlerProxy(assignedEvent)

	originalType := reflect.TypeOf(assignedEvent.handler)
	typedHandler := reflect.MakeFunc(originalType, handlerProxy)

	c.owner.discord.AddHandler(typedHandler.Interface())
}

// === One-Time Handlers

// RegisterOnce registers an event handler as a one-time event handler.
// The registered handler will be removed after being triggered once.
//
//
// The passed handler function will be:
//   1. registered in DiscordGo as a handler
//   2. prepared to allow decorations
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
func (c *ComponentHandlerContainer) RegisterOnce(
	name string,
	handler interface{},
) (string, error) {
	handlerName := GetHandlerName(c.owner, name)

	if _, ok := GetHandler(handlerName); ok {
		return handlerName, fmt.Errorf(
			"an handler for component \"%v\" with name \"%v\" is already registered (ID: \"%v\")",
			c.owner.Name,
			name,
			handlerName)
	}

	assignedEvent := &AssignedEventHandler{
		name:      handlerName,
		component: c.owner,
		handler:   handler,
	}
	c.addComponentHandler(handlerName, assignedEvent)
	err := c.AddDecorator(name, decorateOneTimeHandler)
	if nil != err {
		return handlerName, fmt.Errorf(
			"failed to register cleanup decorator for one-time handler with name \"%v\"",
			handlerName)
	}

	c.addDiscordGoOnceTimeHandler(assignedEvent)

	c.owner.Logger().Info(
		"One-time handler with name \"%v\" for component \"%v\" has been registered! (ID: \"%v\")",
		name,
		c.owner.Name,
		handlerName)

	return handlerName, nil
}

// addDiscordGoOnceTimeHandler generates a handler proxy and registers it for DiscordGo.
func (c *ComponentHandlerContainer) addDiscordGoOnceTimeHandler(assignedEvent *AssignedEventHandler) {
	handlerProxy := createHandlerProxy(assignedEvent)

	originalType := reflect.TypeOf(assignedEvent.handler)
	typedHandler := reflect.MakeFunc(originalType, handlerProxy)

	c.owner.discord.AddHandlerOnce(typedHandler.Interface())
}

// decorateOneTimeHandler is the decorator function used in one-time
// event handlers. It ensures that the executed handler is removed properly.
func decorateOneTimeHandler(
	assignedEvent *AssignedEventHandler,
	session *discordgo.Session,
	event interface{},
	originalHandler interface{},
) {
	removeComponentHandler(assignedEvent.name)

	reflect.ValueOf(originalHandler).Call([]reflect.Value{
		reflect.ValueOf(session),
		reflect.ValueOf(event),
	})
}

// === Handler removal

// Unregister removes the handler with the given name (if existing) from
// the registered handlers.
//
// If the specified handler does not exist, an error will be returned.
func (c *ComponentHandlerContainer) Unregister(name string) error {
	handlerName := GetHandlerName(c.owner, name)
	handler, ok := GetHandler(handlerName)

	if ok {
		return fmt.Errorf(
			"there is no handler called \"%v\" registered that could be unregistered",
			handlerName)
	}

	handler.unregister()
	removeComponentHandler(handlerName)

	return nil
}

// === Handler management

// addComponentHandler adds a new handler to the registered handlers.
//
// Note that adding a handler with a name that is already in the map
// will override the existing handler (but not unregister it from DiscordGo!).
func (c *ComponentHandlerContainer) addComponentHandler(name string, handler *AssignedEventHandler) {
	handlerComponentMapping.Lock()
	defer handlerComponentMapping.Unlock()

	handlerComponentMapping.handlers[name] = handler
}

// removeComponentHandler removes the handler with the specified name
// from the registered handlers.
//
// It also removes all decorators that are assigned to the
// handler.
func removeComponentHandler(name string) {
	handlerComponentMapping.Lock()
	defer handlerComponentMapping.Unlock()

	delete(handlerComponentMapping.handlers, name)
}

// GetHandler returns a handler by its fully qualified name (id).
// The required ID can be obtained using GetHandlerName.
func GetHandler(name string) (*AssignedEventHandler, bool) {
	handlerComponentMapping.RLock()
	defer handlerComponentMapping.RUnlock()

	handler, ok := handlerComponentMapping.handlers[name]

	return handler, ok
}

// GetHandlerName returns the name of a handler for a component.
//
// It acts as the auto-formatter that should be used to retrieve
// handler names.
func GetHandlerName(c *Component, name string) string {
	return util.StringToSnakeCase(fmt.Sprintf("%v_%v", c.Name, name))
}
