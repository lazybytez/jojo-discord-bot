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
	"reflect"
)

// IsComponentEnabled checks if a specific component is currently enabled
// for a specific guild.
func IsComponentEnabled(comp *Component, guildId string) bool {
	if IsCoreComponent(comp) {
		return true
	}

	em := comp.EntityManager()

	regComp, err := em.RegisteredComponent().Get(comp.Code)
	if nil != err {
		comp.Logger().Warn("Missing component with name \"%v\" in database!", comp.Name)
	}

	globalStatus, _ := em.GlobalComponentStatus().Get(regComp.ID)
	if !globalStatus.Enabled {
		return false
	}

	if "" == guildId {
		comp.Logger().Warn("Missing guild with ID \"%v\" in database!", guildId)

		return false
	}
	guild, err := em.Guilds().Get(guildId)
	if nil != err {
		comp.Logger().Warn("Missing guild with ID \"%v\" in database!", comp.Name)
	}

	guildStatus, _ := em.GuildComponentStatus().Get(guild.ID, regComp.ID)

	return guildStatus.Enabled
}

// registerComponentStatusDecorator adds the decorator that handles
// if event handlers should be called or not depending on the owning components
// status
func (c *ComponentHandlerContainer) registerComponentStatusDecorator(name string) {
	ok := c.AddDecorator(name, decorateComponentStatus)
	if !ok {
		c.owner.Logger().Err(fmt.Errorf(
			"failed to register cleanup decorator for one-time Handler with name \"%v\" of component \"%v\"",
			name,
			c.owner.Name),
			"Failed to register one-time handler!")
	}
}

// decorateComponentStatus ensures that handlers are only called when enabled globally
// and enabled on the target guild.
func decorateComponentStatus(
	assignedEvent *AssignedEventHandler,
	session *discordgo.Session,
	event interface{},
	originalHandler interface{},
) {
	comp := assignedEvent.GetComponent()
	guildId := getGuildIdFromEventInterface(event)

	if IsComponentEnabled(comp, guildId) {
		reflect.ValueOf(originalHandler).Call([]reflect.Value{
			reflect.ValueOf(session),
			reflect.ValueOf(event),
		})
	}
}

// getGuildIdFromEventInterface returns the guild id of an event.
// It first tries to get the value of a GuildID field, if this doesn't work,
// the ID field is used (event should be guild event in this case).
//
//	If everything fails, an empty string is returned.
func getGuildIdFromEventInterface(event interface{}) string {
	val := reflect.ValueOf(event).Elem()

	if field := val.FieldByName("GuildID"); field != (reflect.Value{}) {
		return field.String()
	}

	if field := val.FieldByName("ID"); field != (reflect.Value{}) {
		return field.String()
	}

	return ""
}