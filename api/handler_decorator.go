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

// === Structs for decorator management

// decoratorChainElement is an element of the singly-linked
// decoratorChain list.
type decoratorChainElement struct {
    value interface{}
    next  *decoratorChainElement
}

// DecoratorChain is a really simple "singly-linked list" optimized
// for use with the component event handler decorator system.
//
// Using the Obtain function the head element of the chain is returned.
// Then, the next links on the decoratorChainElement can be used to
// walk through the chain.
//
// By using the links on the elements itself, the chain is safe to use concurrently,
// as we do not have a fixed pointer in the chain itself.
type decoratorChain struct {
    head *decoratorChainElement
    tail *decoratorChainElement
}

// Add can be used to add a new element to a decoratorChain.
// The element will be appended at the end of the chain.
func (dC *decoratorChain) Add(decorator interface{}) {
    dCElement := &decoratorChainElement{
        value: decorator,
    }

    if nil == dC.head {
        dC.head = dCElement
        dC.tail = dCElement

        return
    }

    dC.tail.next = dCElement
    dC.tail = dCElement
}

// Obtain returns the first element of the decoratorChain.
func (dC *decoratorChain) Obtain() *decoratorChainElement {
    return dC.head
}

// IsEmpty checks if the chain is empty.
// The chain is considered empty, when no head is present.
func (dC *decoratorChain) IsEmpty() bool {
    return nil == dC.head
}

// Decorator management

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
//
// Note that the name parameter is not the name for the decorator.
// It is the name of the handler that should be decorated
func (c *ComponentHandlerContainer) AddDecorator(name string, decorator interface{}) error {
    handlerName := GetHandlerName(c.owner, name)
    handler, ok := GetHandler(handlerName)

    if !ok {
        return fmt.Errorf(
            "tried to decorate non-existent handler with name \"%v\"",
            handlerName)
    }

    c.appendDecorator(handler, decorator)

    return nil
}

// appendDecorator takes a handler and a decorator and appends it to
// the appropriate decorator list of the AssignedEventHandler
func (c *ComponentHandlerContainer) appendDecorator(handler *AssignedEventHandler, decorator interface{}) {
    if nil == handler.decorators {
        dC := decoratorChain{}

        handler.decorators = &dC
    }

    handler.decorators.Add(decorator)
}
