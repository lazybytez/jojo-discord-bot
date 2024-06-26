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
	"github.com/bwmarrin/discordgo"
	"reflect"
	"sort"
)

// Components holds all available components that
// have been registered.
//
// The component code is used as key.
// The code of a component is unique, registering multiple components
// with the same code results in overriding a previously registered one.
var Components = make([]*Component, 0)

// RegisterComponent adds the given Component to the list
// of registered components.
//
// This function should only being called in the init functions of
// components.
// To get the application to automatically call this function, add an import
// to the <repo-root>/components/component_registry.go.
func RegisterComponent(component *Component, loadComponentFunction func(session *discordgo.Session) error) {
	component.loadComponentFunction = loadComponentFunction
	Components = append(Components, component)

	// This is run once when starting the application.
	// It is not expected to have more than 100 components,
	// therefore just sort after every mutation.
	sortComponents()
}

// sortComponents sorts the components contained in
// the Components slice.
//
// The following logic is applied:
//   - First split components with a code starting with "bot_" out
//   - Sort components with code prefix "bot_" after their priority
//   - Sort other components after their priority
//   - Append sorted normal components to the "bot_" Components
//
// The slice is now sorted after priority, with "bot_" components being always first
func sortComponents() {
	coreComponents := make([]*Component, 0)
	featureComponents := make([]*Component, 0)

	for _, comp := range Components {
		if IsCoreComponent(comp) {
			coreComponents = append(coreComponents, comp)

			continue
		}

		featureComponents = append(featureComponents, comp)
	}

	sort.SliceStable(coreComponents, func(i, j int) bool {
		return coreComponents[i].LoadPriority > coreComponents[j].LoadPriority
	})

	sort.SliceStable(featureComponents, func(i, j int) bool {
		return featureComponents[i].LoadPriority > featureComponents[j].LoadPriority
	})

	Components = []*Component{}
	Components = append(Components, coreComponents...)
	Components = append(Components, featureComponents...)
}

// IsComponentEnabled checks if a specific component is currently enabled
// for a specific guild.
// If the guild id is empty, the function will return the global status of the component.
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
		return true
	}

	guild, err := em.Guilds().Get(guildId)
	if nil != err {
		comp.Logger().Warn("Missing guild with ID \"%v\" in database!", comp.Name)
	}

	guildStatus, _ := em.GuildComponentStatus().Get(guild.ID, regComp.ID)

	return guildStatus.Enabled
}

// getGuildIdFromEventInterface returns the guild id of an event.
// It first tries to get the value of a GuildID field, if this doesn't work,
// the ID field is used (event should be guild event in this case).
//
//	If everything fails, an empty string is returned.
func getGuildIdFromEventInterface(event interface{}) string {
	val := reflect.ValueOf(event)
	if reflect.Pointer == val.Kind() {
		val = val.Elem()
	}

	if field := val.FieldByName("GuildID"); field != (reflect.Value{}) {
		return field.String()
	}

	if field := val.FieldByName("ID"); field != (reflect.Value{}) {
		return field.String()
	}

	return ""
}
