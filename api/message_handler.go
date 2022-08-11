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
	"regexp"
)

// CompoMessageHandlerManager is an interface that provides functions to register specialized
// event handlers for message handling.
type CompoMessageHandlerManager interface {
	// RegisterSimpleMessageHandler can be used to register a new AssignedEventHandler to handle messages.
	//
	// In general, the function just works like the common Register method of the API,
	// except that strings can be passed as third argument.
	// The created event handler will be decorated and only trigger if the received message
	// matches one of the specified strings.
	//
	// Note that handlers registered using this function ignore messages from bots
	// and messages that the bot itself send.
	//
	// The passed Handler function will be:
	//  1. registered in DiscordGo as a Handler
	//  2. prepared to allow decorations
	//  3. saved with a name that allows to retrieve it later on
	//
	// The passed name for the Handler is concatenated with the name of the
	// component that owns the Handler (separated by underscore).
	//
	// The Handler must have the same format as when a Handler is registered in
	// plain DiscordGo. See the documentation about discordgo.AddHandler
	// for additional information.
	//
	// In general, the common format for a Handler function is:
	//
	//	func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
	RegisterSimpleMessageHandler(name string, handler interface{}, messages ...string) (string, bool)

	// RegisterComplexMessageHandler can be used to register a new AssignedEventHandler to handle messages.
	//
	// In general, the function just works like the common Register method of the API,
	// except that strings can be passed as third argument.
	// The created event handler will be decorated and only trigger if the received message
	// matches one of the specified strings.
	//
	// The passed messages should be regex patterns. What makes RegisterComplexMessageHandler special
	// is the fact that it allows to match messages with Regex.
	//
	// Note that handlers registered using this function ignore messages from bots
	// and messages that the bot itself send.
	//
	// The passed Handler function will be:
	//  1. registered in DiscordGo as a Handler
	//  2. prepared to allow decorations
	//  3. saved with a name that allows to retrieve it later on
	//
	// The passed name for the Handler is concatenated with the name of the
	// component that owns the Handler (separated by underscore).
	//
	// The Handler must have the same format as when a Handler is registered in
	// plain DiscordGo. See the documentation about discordgo.AddHandler
	// for additional information.
	//
	// In general, the common format for a Handler function is:
	//
	//	func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
	//
	// The returned values contain the name of the (potentially registered) handler
	// and a bool indicating whether the registration was successful or not.
	RegisterComplexMessageHandler(name string, handler interface{}, messages ...string) (string, bool)
}

// RegisterSimpleMessageHandler can be used to register a new AssignedEventHandler to handle messages.
//
// In general, the function just works like the common Register method of the API,
// except that strings can be passed as third argument.
// The created event handler will be decorated and only trigger if the received message
// matches one of the specified strings.
//
// Note that handlers registered using this function ignore messages from bots
// and messages that the bot itself send.
//
// The passed Handler function will be:
//  1. registered in DiscordGo as a Handler
//  2. prepared to allow decorations
//  3. saved with a name that allows to retrieve it later on
//
// The passed name for the Handler is concatenated with the name of the
// component that owns the Handler (separated by underscore).
//
// The Handler must have the same format as when a Handler is registered in
// plain DiscordGo. See the documentation about discordgo.AddHandler
// for additional information.
//
// In general, the common format for a Handler function is:
//
//	func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
//
// The returned values contain the name of the (potentially registered) handler
// and a bool indicating whether the registration was successful or not.
func (c *ComponentHandlerContainer) RegisterSimpleMessageHandler(
	name string,
	handler interface{},
	messages ...string,
) (string, bool) {
	if nil == messages || len(messages) < 1 {
		c.owner.Logger().Err(fmt.Errorf(
			"at least one message is expected for a simple message handler, "+
				"<=0 messages have been passed for handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register simple message handler")

		return name, false
	}

	handlerName, ok := c.Register(name, handler)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to register simple message handler for handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register simple message handler")

		return handlerName, false
	}

	ok = c.AddDecorator(name, decorateSimpleMessageHandler)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to apply necessary simple message handler decorator for, "+
				"handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register simple message handler")

		return handlerName, false
	}

	// Until this point original handler would be registered not being called
	// due to missing messages. After this point setup is complete
	assignedHandler, ok := GetHandler(handlerName)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to get assigned handler for handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name), "Failed to register simple message handler")

		return handlerName, false
	}
	assignedHandler.handledMessages = messages

	return handlerName, true
}

// decorateSimpleMessageHandler is a decorator that ensures a handler is only called
// when a received message originates from a user and matches a specific string.
func decorateSimpleMessageHandler(
	assignedEvent *AssignedEventHandler,
	session *discordgo.Session,
	event *discordgo.MessageCreate,
	originalHandler func(s *discordgo.Session, m *discordgo.MessageCreate),
) {
	if event.Author.ID == session.State.User.ID || event.Author.Bot {
		return
	}

	content := event.Content
	for _, message := range assignedEvent.handledMessages {
		if message == content {
			originalHandler(session, event)

			return
		}
	}

	return
}

// RegisterComplexMessageHandler can be used to register a new AssignedEventHandler to handle messages.
//
// In general, the function just works like the common Register method of the API,
// except that strings can be passed as third argument.
// The created event handler will be decorated and only trigger if the received message
// matches one of the specified strings.
//
// The passed messages should be regex patterns. What makes RegisterComplexMessageHandler special
// is the fact that it allows to match messages with Regex.
//
// Note that handlers registered using this function ignore messages from bots
// and messages that the bot itself send.
//
// The passed Handler function will be:
//  1. registered in DiscordGo as a Handler
//  2. prepared to allow decorations
//  3. saved with a name that allows to retrieve it later on
//
// The passed name for the Handler is concatenated with the name of the
// component that owns the Handler (separated by underscore).
//
// The Handler must have the same format as when a Handler is registered in
// plain DiscordGo. See the documentation about discordgo.AddHandler
// for additional information.
//
// In general, the common format for a Handler function is:
//
//	func (session *discordgo.Session, event <event to call, e.g. discordgo.MessageCreate)
//
// The returned values contain the name of the (potentially registered) handler
// and a bool indicating whether the registration was successful or not.
func (c *ComponentHandlerContainer) RegisterComplexMessageHandler(
	name string,
	handler interface{},
	patterns ...string,
) (string, bool) {
	if nil == patterns || len(patterns) < 1 {
		c.owner.Logger().Err(fmt.Errorf(
			"at least one regex is expected for a complex message handler, "+
				"<=0 patterns have been passed for handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register complex message handler")

		return name, false
	}

	handlerName, ok := c.Register(name, handler)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to register complex message handler for handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register complex message handler")

		return handlerName, false
	}

	ok = c.AddDecorator(name, decorateComplexMessageHandler)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to apply necessary complex message handler decorator for, "+
				"handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register complex message handler")

		return handlerName, false
	}

	// Until this point original handler would be registered not being called
	// due to missing patterns. After this point setup is complete
	assignedHandler, ok := GetHandler(handlerName)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to get assigned handler for handler \"%v\" of component \"%v\"",
			name,
			c.owner.Name), "Failed to register complex message handler")

		return handlerName, false
	}

	compiledPatterns := make([]*regexp.Regexp, 0)
	for _, pattern := range patterns {
		compiledPatterns = append(compiledPatterns, regexp.MustCompile(pattern))
	}
	assignedHandler.handledPatterns = compiledPatterns

	return handlerName, true
}

// decorateComplexMessageHandler is a decorator that ensures a handler is only called
// when a received message originates from a user and matches a specific regex.
func decorateComplexMessageHandler(
	assignedEvent *AssignedEventHandler,
	session *discordgo.Session,
	event *discordgo.MessageCreate,
	originalHandler func(s *discordgo.Session, m *discordgo.MessageCreate),
) {
	if event.Author.ID == session.State.User.ID || event.Author.Bot {
		return
	}

	content := event.Content
	for _, pattern := range assignedEvent.handledPatterns {
		if pattern.MatchString(content) {
			originalHandler(session, event)

			return
		}
	}

	return
}
