package api

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/lazybytez/jojo-discord-bot/api/util"
    "golang.org/x/crypto/openpgp/errors"
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

// componentHandlerMap is a wrapper that holds the map that contains
// the handler name -> handler mapping. It embeds sync.RWMutex
// to allow support for concurrency.
type componentHandlerMap struct {
    sync.RWMutex
    handlers map[string]AssignedEventHandler
    //decorators map[string][]string
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
    handlers: make(map[string]AssignedEventHandler),
}

// AssignedEventHandler holds the necessary data to handle events
// in the bot. It is used to allow features like decorating handlers
// by caching the original handler function.
// The original handler function is then replaced with a proxy function
// that allows us to add special features.
type AssignedEventHandler struct {
    name      string
    component *Component
    handler   interface{}

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
    Register(name string, handler interface{}) (string, error)
    RegisterOnce(name string, handler interface{})
    Unregister(name string)

    addDiscordGoHandler(assignedEvent AssignedEventHandler)
    addComponentHandler(name string, handler AssignedEventHandler)
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
        c.handlerManager = ComponentHandlerContainer{
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
// The handler must have the same format as when a handler is registered in
// plain DiscordGo. See the documentation about discordgo.AddHandler
// for additional information.
//
// In general, the common format for a handler function is:
//   func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
func (c ComponentHandlerContainer) Register(name string, handler interface{}) (string, error) {
    handlerName := GetHandlerName(c.owner, name)

    if _, ok := GetHandler(handlerName); ok {
        return handlerName, errors.InvalidArgumentError(fmt.Sprintf(
            "An handler for component \"%v\" with name \"%v\" is already registered (ID: \"%v\")",
            c.owner.Name,
            name,
            handlerName))
    }

    assignedEvent := AssignedEventHandler{
        name:      handlerName,
        component: c.owner,
        handler:   handler,
    }
    c.addComponentHandler(handlerName, assignedEvent)

    c.addDiscordGoHandler(assignedEvent)

    return handlerName, nil
}

// createHandlerProxy creates a closure that acts as a proxy for the original handler function.
// This allows us to delegate the triggered event to the decorateHandler function.
// The decorateHandler than allows to intercept the original event handler with
// custom decorator functions. This can be useful in some edge-cases.
func createHandlerProxy(handler AssignedEventHandler) func(args []reflect.Value) []reflect.Value {
    return func(args []reflect.Value) []reflect.Value {
        decorateArguments := append([]reflect.Value{reflect.ValueOf(handler)}, args...)
        reflect.ValueOf(decorateHandler).Call(decorateArguments)

        return []reflect.Value{}
    }
}

// decorateHandler is the internal replacement used when an event happens in discord.
func decorateHandler(assignedEvent AssignedEventHandler, session *discordgo.Session, e interface{}) {
    if _, ok := GetHandler(assignedEvent.name); !ok {
        assignedEvent.component.Logger().Warn(fmt.Sprintf(
            "Potentially orphaned event handler named \"%v\" has been called! "+
                "Ensure to properly unregister no longer needed handlers!",
            assignedEvent.name))

        return
    }

    sessionRef := reflect.ValueOf(session)
    eRef := reflect.ValueOf(e)

    reflect.ValueOf(assignedEvent.handler).Call([]reflect.Value{sessionRef, eRef})
}

// addDiscordGoHandler generates a handler proxy and registers it for DiscordGo.
// This function first
func (c ComponentHandlerContainer) addDiscordGoHandler(assignedEvent AssignedEventHandler) {
    handlerProxy := createHandlerProxy(assignedEvent)

    originalType := reflect.TypeOf(assignedEvent.handler)
    typedHandler := reflect.MakeFunc(originalType, handlerProxy)

    c.owner.discord.AddHandler(typedHandler.Interface())
}

// addComponentHandler adds a new handler to the registered handlers.
//
// Note that adding a handler with a name that is already in the map
// will override the existing handler (but not unregister it from DiscordGo!).
func (c ComponentHandlerContainer) addComponentHandler(name string, handler AssignedEventHandler) {
    handlerComponentMapping.Lock()
    defer handlerComponentMapping.Unlock()

    handlerComponentMapping.handlers[name] = handler
}

// GetHandler returns a handler by its fully qualified name (id).
// The required ID can be obtained using GetHandlerName.
func GetHandler(name string) (AssignedEventHandler, bool) {
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

// RegisterOnce TODO: Implement when decorating is possible
func (c ComponentHandlerContainer) RegisterOnce(name string, handler interface{}) {
}

func (c ComponentHandlerContainer) Unregister(name string) {

}
