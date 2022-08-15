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

package bot_core

import (
	"github.com/lazybytez/jojo-discord-bot/api"
)

// registerAvailableComponents ensures that all available components
// are registered in the database.
//
// Note that this function only adds new components, we do not care
// about orphaned components.
func registerAvailableComponents() {
	for _, component := range api.Components { // safe to assume that Components is populated
		registeredComponent, err := C.EntityManager().RegisteredComponent().Get(component.Code)

		if nil != err {
			registeredComponent.Code = component.Code

			err := C.EntityManager().Create(registeredComponent)
			if nil != err {
				C.Logger().Warn(
					"Failed to register component with code \"%v\" in database!",
					registeredComponent.Code)
			}
		}
	}
}

// ensureGlobalComponentStatusExists ensures that for every component
// a database entry in the global status table exists.
//
// By default, this is created with an enabled status, as the global status
// is only meant for maintenance purposes.
func ensureGlobalComponentStatusExists() {
	for _, component := range api.Components { // safe to assume that Components is populated
		registeredComponent, err := C.EntityManager().RegisteredComponent().Get(component.Code)

		if nil != err {
			continue
		}

		globalComponentStatus, err := C.EntityManager().GlobalComponentStatus().Get(registeredComponent.ID)
		if err != nil {
			globalComponentStatus.Component = *registeredComponent
			globalComponentStatus.Enabled = true

			err := C.EntityManager().Create(globalComponentStatus)
			if err != nil {
				C.Logger().Warn(
					"Failed to global component status for component \"%v\" in database!",
					registeredComponent.Code)
			}
		}
	}
}
