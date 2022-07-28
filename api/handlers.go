package api

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/lazybytez/jojo-discord-bot/api/util"
    "golang.org/x/crypto/openpgp/errors"
    "reflect"
    "sync"
)

type DiscordGoHandler func(session *discordgo.Session, e interface{})

// componentHandlerMap is a wrapper that holds the map that contains
// the handler name -> handler mapping. It embeds sync.RWMutex
// to allow support for concurrency.
type componentHandlerMap struct {
    sync.RWMutex
    handlers map[string]AssignedEventHandler
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
    GetHandler() DiscordGoHandler
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
    RegisterOnce(name string, handler DiscordGoHandler)
    Unregister(name string)
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

    handlerProxy := func(session *discordgo.Session, e interface{}) {
        decorateHandler(assignedEvent, session, e)
    }
    assignedEvent.handler = handler
    addHandler(handlerName, assignedEvent)
    c.owner.discord.AddHandler(handlerProxy)

    return handlerName, nil
}

// addHandler adds a new handler to the registered handlers.
//
// Note that adding a handler with a name that is already in the map
// will override the existing handler (but not unregister it from DiscordGo!).
func addHandler(name string, handler AssignedEventHandler) {
    handlerComponentMapping.Lock()
    handlerComponentMapping.handlers[name] = handler
    handlerComponentMapping.Unlock()
}

// GetHandler returns a handler by its fully qualified name (id).
// The required ID can be obtained using GetHandlerName.
func GetHandler(name string) (AssignedEventHandler, bool) {
    handlerComponentMapping.RLock()
    handler, ok := handlerComponentMapping.handlers[name]
    handlerComponentMapping.RUnlock()

    return handler, ok
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

    fmt.Println(reflect.TypeOf(assignedEvent.handler))
    switch c := assignedEvent.handler.(type) {
    case func(*discordgo.Session, discordgo.Event):
        c(session, e.(discordgo.Event))
    default:
        assignedEvent.component.Logger().Warn(fmt.Sprintf(
            "Failed to call handler with id \"%v\"! The type of the signature in the passed handler is wrong!",
            assignedEvent.name))
    }
}

func (c ComponentHandlerContainer) RegisterOnce(name string, handler DiscordGoHandler) {
}

func (c ComponentHandlerContainer) Unregister(name string) {

}

// GetHandlerName returns the name of a handler for a component.
//
// It acts as the auto-formatter that should be used to retrieve
// handler names.
func GetHandlerName(c *Component, name string) string {
    return util.StringToSnakeCase(fmt.Sprintf("%v_%v", c.Name, name))
}
